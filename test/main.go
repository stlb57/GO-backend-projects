package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

const (
	expandedWidth  = 340
	collapsedWidth = 18
	windowHeight   = 620

	maxRecent  = 12
	maxResults = 100
)

type Folder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	Permanent []Folder `json:"permanent"`
	Recent    []Folder `json:"recent"`
}

type FileResult struct {
	Name    string
	Path    string
	Size    int64
	ModTime time.Time
}

var (
	config     Config
	configPath string

	currentFolder string

	sessionPins []Folder

	w                  fyne.Window
	panel              *fyne.Container
	handle             *widget.Button
	permanentContainer *fyne.Container
	sessionContainer   *fyne.Container
	recentContainer    *fyne.Container
	resultContainer    *fyne.Container

	currentFolderLabel *widget.Label
	statusLabel        *widget.Label
	searchEntry        *widget.Entry
	sortSelect         *widget.Select

	uiMu sync.Mutex

	user32                  = syscall.NewLazyDLL("user32.dll")
	getForegroundWindowProc = user32.NewProc("GetForegroundWindow")
	getSystemMetricsProc    = user32.NewProc("GetSystemMetrics")
	getCursorPosProc        = user32.NewProc("GetCursorPos")
	getWindowLongPtrProc    = user32.NewProc("GetWindowLongPtrW")
	setWindowLongPtrProc    = user32.NewProc("SetWindowLongPtrW")
	setLayeredProc          = user32.NewProc("SetLayeredWindowAttributes")
	setWindowPosProc        = user32.NewProc("SetWindowPos")
)

type point struct {
	X int32
	Y int32
}

func main() {
	a := app.NewWithID("com.filevault.app")
	a.Settings().SetTheme(theme.DarkTheme())

	w = a.NewWindow("FileVault")
	w.SetPadded(false)
	w.SetFixedSize(true)

	loadConfig()

	buildUI()

	w.Resize(fyne.NewSize(
		collapsedWidth,
		windowHeight,
	))

	positionWindow(collapsedWidth)

	if desk, ok := w.(desktop.Window); ok {
		desk.RequestAlwaysOnTop()
	}

	startExplorerTracker()
	startEdgeWatcher()

	w.SetCloseIntercept(func() {
		saveConfig()
		a.Quit()
	})

	w.Show()

	applyWindowsWindowStyle()

	a.Run()
}

func buildUI() {
	handle = widget.NewButton("‹", func() {
		expandPanel()
	})

	handle.Alignment = widget.ButtonAlignCenter

	currentFolderLabel = widget.NewLabel(
		"No Explorer folder detected",
	)

	currentFolderLabel.Wrapping = fyne.TextWrapWord

	statusLabel = widget.NewLabel("Ready")

	searchEntry = widget.NewEntry()
	searchEntry.SetPlaceHolder("Search current folder...")

	searchEntry.OnSubmitted = func(query string) {
		searchFiles(query)
	}

	sortSelect = widget.NewSelect(
		[]string{
			"Newest",
			"Oldest",
			"Largest",
			"Name",
		},
		func(string) {
			if searchEntry.Text != "" {
				searchFiles(searchEntry.Text)
			}
		},
	)

	sortSelect.SetSelected("Newest")

	searchButton := widget.NewButton("Search", func() {
		searchFiles(searchEntry.Text)
	})

	searchRow := container.NewBorder(
		nil,
		nil,
		nil,
		searchButton,
		searchEntry,
	)

	permanentContainer = container.NewVBox()
	sessionContainer = container.NewVBox()
	recentContainer = container.NewVBox()
	resultContainer = container.NewVBox()

	pinCurrentButton := widget.NewButton(
		"⚡ Pin current",
		func() {
			if currentFolder == "" {
				statusLabel.SetText(
					"No Explorer folder detected.",
				)
				return
			}

			addSessionPin(currentFolder)
		},
	)

	addPermanentButton := widget.NewButton(
		"+ Permanent",
		func() {
			if currentFolder == "" {
				statusLabel.SetText(
					"No Explorer folder detected.",
				)
				return
			}

			addPermanent(currentFolder)
		},
	)

	sortRow := container.NewBorder(
		nil,
		nil,
		widget.NewLabel("Sort"),
		nil,
		sortSelect,
	)

	panel = container.NewVBox(
		widget.NewLabelWithStyle(
			"FILEVAULT",
			fyne.TextAlignLeading,
			fyne.TextStyle{
				Bold: true,
			},
		),

		currentFolderLabel,

		widget.NewSeparator(),

		widget.NewLabel("📌  PERMANENT"),
		permanentContainer,
		addPermanentButton,

		widget.NewSeparator(),

		widget.NewLabel("⚡  SESSION"),
		sessionContainer,
		pinCurrentButton,

		widget.NewSeparator(),

		widget.NewLabel("🕘  RECENT"),
		recentContainer,

		widget.NewSeparator(),

		widget.NewLabel("🔎  FILE SEARCH"),
		searchRow,
		sortRow,
		resultContainer,

		widget.NewSeparator(),

		statusLabel,
	)

	content := container.NewMax(
		handle,
		panel,
	)

	panel.Hide()

	w.SetContent(content)

	refreshPermanentUI()
	refreshSessionUI()
	refreshRecentUI()
}

func refreshPermanentUI() {
	permanentContainer.RemoveAll()

	for _, folder := range config.Permanent {
		folder := folder

		openButton := widget.NewButton(
			folder.Name,
			func() {
				openFolder(folder.Path)
			},
		)

		removeButton := widget.NewButton(
			"×",
			func() {
				removePermanent(folder.Path)
			},
		)

		row := container.NewBorder(
			nil,
			nil,
			nil,
			removeButton,
			openButton,
		)

		permanentContainer.Add(row)
	}

	permanentContainer.Refresh()
}

func refreshSessionUI() {
	sessionContainer.RemoveAll()

	for _, folder := range sessionPins {
		folder := folder

		openButton := widget.NewButton(
			folder.Name,
			func() {
				openFolder(folder.Path)
			},
		)

		removeButton := widget.NewButton(
			"×",
			func() {
				removeSessionPin(folder.Path)
			},
		)

		row := container.NewBorder(
			nil,
			nil,
			nil,
			removeButton,
			openButton,
		)

		sessionContainer.Add(row)
	}

	if len(sessionPins) == 0 {
		sessionContainer.Add(
			widget.NewLabel(
				"Nothing pinned for this session.",
			),
		)
	}

	sessionContainer.Refresh()
}

func refreshRecentUI() {
	recentContainer.RemoveAll()

	for _, folder := range config.Recent {
		folder := folder

		openButton := widget.NewButton(
			folder.Name,
			func() {
				openFolder(folder.Path)
			},
		)

		pinButton := widget.NewButton(
			"⚡",
			func() {
				addSessionPin(folder.Path)
			},
		)

		row := container.NewBorder(
			nil,
			nil,
			nil,
			pinButton,
			openButton,
		)

		recentContainer.Add(row)
	}

	if len(config.Recent) == 0 {
		recentContainer.Add(
			widget.NewLabel(
				"Open a folder in Explorer...",
			),
		)
	}

	recentContainer.Refresh()
}

func addPermanent(path string) {
	path = filepath.Clean(path)

	for _, folder := range config.Permanent {
		if strings.EqualFold(folder.Path, path) {
			statusLabel.SetText(
				"Already permanently pinned.",
			)
			return
		}
	}

	config.Permanent = append(
		config.Permanent,
		Folder{
			Name: displayName(path),
			Path: path,
		},
	)

	saveConfig()
	refreshPermanentUI()

	statusLabel.SetText(
		"Pinned permanently.",
	)
}

func removePermanent(path string) {
	newList := make(
		[]Folder,
		0,
		len(config.Permanent),
	)

	for _, folder := range config.Permanent {
		if !strings.EqualFold(
			folder.Path,
			path,
		) {
			newList = append(
				newList,
				folder,
			)
		}
	}

	config.Permanent = newList

	saveConfig()
	refreshPermanentUI()
}

func addSessionPin(path string) {
	path = filepath.Clean(path)

	for _, folder := range sessionPins {
		if strings.EqualFold(
			folder.Path,
			path,
		) {
			statusLabel.SetText(
				"Already pinned for this session.",
			)
			return
		}
	}

	sessionPins = append(
		sessionPins,
		Folder{
			Name: displayName(path),
			Path: path,
		},
	)

	refreshSessionUI()

	statusLabel.SetText(
		"Pinned for this session.",
	)
}

func removeSessionPin(path string) {
	newList := make(
		[]Folder,
		0,
		len(sessionPins),
	)

	for _, folder := range sessionPins {
		if !strings.EqualFold(
			folder.Path,
			path,
		) {
			newList = append(
				newList,
				folder,
			)
		}
	}

	sessionPins = newList

	refreshSessionUI()
}

func addRecent(path string) {
	path = filepath.Clean(path)

	if path == "" {
		return
	}

	for i, folder := range config.Recent {
		if strings.EqualFold(
			folder.Path,
			path,
		) {
			if i != 0 {
				config.Recent = append(
					[]Folder{folder},
					append(
						config.Recent[:i],
						config.Recent[i+1:]...,
					)...,
				)
			}

			saveConfig()
			refreshRecentUI()

			return
		}
	}

	config.Recent = append(
		[]Folder{
			{
				Name: displayName(path),
				Path: path,
			},
		},
		config.Recent...,
	)

	if len(config.Recent) > maxRecent {
		config.Recent =
			config.Recent[:maxRecent]
	}

	saveConfig()
	refreshRecentUI()
}

func openFolder(path string) {
	info, err := os.Stat(path)

	if err != nil || !info.IsDir() {
		statusLabel.SetText(
			"Folder no longer exists.",
		)
		return
	}

	cmd := exec.Command(
		"explorer.exe",
		"/n,",
		path,
	)

	if err := cmd.Start(); err != nil {
		statusLabel.SetText(
			"Could not open Explorer.",
		)
		return
	}

	statusLabel.SetText(
		"Opening " + displayName(path),
	)
}

func searchFiles(query string) {
	query = strings.TrimSpace(query)

	if query == "" {
		resultContainer.RemoveAll()
		resultContainer.Refresh()
		return
	}

	if currentFolder == "" {
		statusLabel.SetText(
			"No Explorer folder detected.",
		)
		return
	}

	root := currentFolder

	resultContainer.RemoveAll()

	resultContainer.Add(
		widget.NewLabel("Searching..."),
	)

	resultContainer.Refresh()

	go func() {
		results := scanFiles(
			root,
			query,
		)

		sortResults(results)

		if len(results) > maxResults {
			results = results[:maxResults]
		}

		fyne.Do(func() {
			resultContainer.RemoveAll()

			for _, result := range results {
				result := result

				button := widget.NewButton(
					result.Name+
						"  •  "+
						formatSize(result.Size),
					func() {
						openFolder(
							filepath.Dir(
								result.Path,
							),
						)
					},
				)

				resultContainer.Add(button)
			}

			if len(results) == 0 {
				resultContainer.Add(
					widget.NewLabel(
						"No matching files.",
					),
				)
			}

			resultContainer.Refresh()

			statusLabel.SetText(
				strconv.Itoa(len(results)) +
					" result(s)",
			)
		})
	}()
}

func scanFiles(
	root string,
	query string,
) []FileResult {
	var results []FileResult

	pattern := normalizePattern(query)

	filepath.WalkDir(
		root,
		func(
			path string,
			entry os.DirEntry,
			err error,
		) error {
			if err != nil {
				return nil
			}

			if entry.IsDir() {
				return nil
			}

			name := entry.Name()

			matched, err := filepath.Match(
				pattern,
				strings.ToLower(name),
			)

			if err != nil {
				return nil
			}

			if !matched {
				return nil
			}

			info, err := entry.Info()

			if err != nil {
				return nil
			}

			results = append(
				results,
				FileResult{
					Name:    name,
					Path:    path,
					Size:    info.Size(),
					ModTime: info.ModTime(),
				},
			)

			return nil
		},
	)

	return results
}

func normalizePattern(query string) string {
	query = strings.ToLower(
		strings.TrimSpace(query),
	)

	query = strings.ReplaceAll(
		query,
		"%",
		"*",
	)

	query = strings.ReplaceAll(
		query,
		"_",
		"?",
	)

	if !strings.ContainsAny(
		query,
		"*?",
	) {
		query = "*" + query + "*"
	}

	return query
}

func sortResults(results []FileResult) {
	switch sortSelect.Selected {
	case "Oldest":
		sort.Slice(
			results,
			func(i, j int) bool {
				return results[i].ModTime.Before(
					results[j].ModTime,
				)
			},
		)

	case "Largest":
		sort.Slice(
			results,
			func(i, j int) bool {
				return results[i].Size >
					results[j].Size
			},
		)

	case "Name":
		sort.Slice(
			results,
			func(i, j int) bool {
				return strings.ToLower(
					results[i].Name,
				) <
					strings.ToLower(
						results[j].Name,
					)
			},
		)

	default:
		sort.Slice(
			results,
			func(i, j int) bool {
				return results[i].ModTime.After(
					results[j].ModTime,
				)
			},
		)
	}
}

func startExplorerTracker() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		err := ole.CoInitializeEx(
			0,
			ole.COINIT_APARTMENTTHREADED,
		)

		if err != nil {
			return
		}

		defer ole.CoUninitialize()

		for {
			path := activeExplorerFolder()

			if path != "" {
				fyne.Do(func() {
					if !strings.EqualFold(
						currentFolder,
						path,
					) {
						currentFolder = path

						currentFolderLabel.SetText(
							"Current\n" + path,
						)

						addRecent(path)
					}
				})
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func activeExplorerFolder() string {
	foreground, _, _ :=
		getForegroundWindowProc.Call()

	if foreground == 0 {
		return ""
	}

	shellUnknown, err :=
		oleutil.CreateObject(
			"Shell.Application",
		)

	if err != nil {
		return ""
	}

	shell, err :=
		shellUnknown.QueryInterface(
			ole.IID_IDispatch,
		)

	shellUnknown.Release()

	if err != nil {
		return ""
	}

	defer shell.Release()

	windowsVariant, err :=
		oleutil.GetProperty(
			shell,
			"Windows",
		)

	if err != nil {
		return ""
	}

	defer windowsVariant.Clear()

	windows :=
		windowsVariant.ToIDispatch()

	if windows == nil {
		return ""
	}

	defer windows.Release()

	var result string

	_ = oleutil.ForEach(
		windows,
		func(
			itemVariant *ole.VARIANT,
		) error {
			if result != "" {
				return nil
			}

			item :=
				itemVariant.ToIDispatch()

			if item == nil {
				return nil
			}

			defer item.Release()

			hwndVariant, err :=
				oleutil.GetProperty(
					item,
					"HWND",
				)

			if err != nil {
				return nil
			}

			hwnd := variantUintptr(
				hwndVariant,
			)

			hwndVariant.Clear()

			if hwnd != foreground {
				return nil
			}

			fullNameVariant, err :=
				oleutil.GetProperty(
					item,
					"FullName",
				)

			if err != nil {
				return nil
			}

			fullName :=
				fullNameVariant.ToString()

			fullNameVariant.Clear()

			if !strings.EqualFold(
				filepath.Base(fullName),
				"explorer.exe",
			) {
				return nil
			}

			locationVariant, err :=
				oleutil.GetProperty(
					item,
					"LocationURL",
				)

			if err != nil {
				return nil
			}

			location :=
				locationVariant.ToString()

			locationVariant.Clear()

			result =
				fileURLToPath(location)

			return nil
		},
	)

	return result
}

func variantUintptr(
	variant *ole.VARIANT,
) uintptr {
	value := variant.Value()

	switch value := value.(type) {
	case int:
		return uintptr(value)

	case int16:
		return uintptr(value)

	case int32:
		return uintptr(uint32(value))

	case int64:
		return uintptr(value)

	case uint:
		return uintptr(value)

	case uint16:
		return uintptr(value)

	case uint32:
		return uintptr(value)

	case uint64:
		return uintptr(value)

	case string:
		number, err :=
			strconv.ParseUint(
				value,
				10,
				64,
			)

		if err == nil {
			return uintptr(number)
		}
	}

	return 0
}

func fileURLToPath(
	value string,
) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	if !strings.HasPrefix(
		strings.ToLower(value),
		"file:",
	) {
		return value
	}

	parsed, err :=
		url.Parse(value)

	if err != nil {
		return ""
	}

	path, err :=
		url.PathUnescape(
			parsed.Path,
		)

	if err != nil {
		return ""
	}

	if parsed.Host != "" {
		path = `\\` + parsed.Host + path
	} else if len(path) >= 3 &&
		path[0] == '/' &&
		((path[1] >= 'a' && path[1] <= 'z') ||
			(path[1] >= 'A' && path[1] <= 'Z')) &&
		path[2] == ':' {
		path = path[1:]
	}

	return filepath.Clean(filepath.FromSlash(path))
}

func startEdgeWatcher() {
	go func() {
		for {
			x, _ := cursorPosition()

			screenWidth :=
				getScreenWidth()

			if panel.Visible() {
				if x <
					screenWidth-
						expandedWidth-
						12 {
					collapsePanel()
				}
			} else {
				if x >=
					screenWidth-
						collapsedWidth-
						10 {
					expandPanel()
				}
			}

			time.Sleep(
				50 * time.Millisecond,
			)
		}
	}()
}

func expandPanel() {
	if panel.Visible() {
		return
	}

	panel.Show()
	handle.Hide()

	w.Resize(
		fyne.NewSize(
			expandedWidth,
			windowHeight,
		),
	)

	positionWindow(
		expandedWidth,
	)

	panel.Refresh()
}

func collapsePanel() {
	if !panel.Visible() {
		return
	}

	panel.Hide()
	handle.Show()

	w.Resize(
		fyne.NewSize(
			collapsedWidth,
			windowHeight,
		),
	)

	positionWindow(
		collapsedWidth,
	)
}

func cursorPosition() (int, int) {
	var p point

	ret, _, _ :=
		getCursorPosProc.Call(
			uintptr(
				unsafe.Pointer(&p),
			),
		)

	if ret == 0 {
		return 0, 0
	}

	return int(p.X), int(p.Y)
}

func getScreenWidth() int {
	ret, _, _ :=
		getSystemMetricsProc.Call(0)

	return int(ret)
}

func positionWindow(width int) {
	screenWidth :=
		getScreenWidth()

	x := screenWidth - width

	if desk, ok :=
		w.(desktop.Window); ok {
		desk.RequestPosition(
			x,
			180,
		)
	}
}

func applyWindowsWindowStyle() {
	nativeWindow, ok :=
		w.(driver.NativeWindow)

	if !ok {
		return
	}

	indexStyle := int32(-16)
	indexExStyle := int32(-20)

	nativeWindow.RunNative(
		func(context any) {
			windowsContext, ok :=
				context.(driver.WindowsWindowContext)

			if !ok {
				return
			}

			hwnd :=
				windowsContext.HWND

			// Remove normal window decorations.
			style, _, _ :=
				getWindowLongPtrProc.Call(
					hwnd,
					uintptr(indexStyle),
				)

			style &^=
				0x00C00000 |
					0x00040000 |
					0x00020000 |
					0x00010000 |
					0x00080000

			setWindowLongPtrProc.Call(
				hwnd,
				uintptr(indexStyle),
				style,
			)

			// Tool window + layered window.
			exStyle, _, _ :=
				getWindowLongPtrProc.Call(
					hwnd,
					uintptr(indexExStyle),
				)

			exStyle |=
				0x00000080 |
					0x00080000

			setWindowLongPtrProc.Call(
				hwnd,
				uintptr(indexExStyle),
				exStyle,
			)

			// 235 / 255 opacity.
			setLayeredProc.Call(
				hwnd,
				0,
				235,
				0x00000002,
			)

			// Apply the style changes.
			setWindowPosProc.Call(
				hwnd,
				^uintptr(0),
				0,
				0,
				0,
				0,
				(0x0001 |
					0x0002 |
					0x0020 |
					0x0010),
			)
		},
	)
}

func loadConfig() {
	dir, err :=
		os.UserConfigDir()

	if err != nil {
		return
	}

	configDir :=
		filepath.Join(
			dir,
			"FileVault",
		)

	if err := os.MkdirAll(
		configDir,
		0755,
	); err != nil {
		return
	}

	configPath =
		filepath.Join(
			configDir,
			"config.json",
		)

	data, err :=
		os.ReadFile(configPath)

	if err != nil {
		return
	}

	_ = json.Unmarshal(
		data,
		&config,
	)

	cleanConfig()
}

func cleanConfig() {
	validPermanent :=
		make([]Folder, 0)

	for _, folder := range config.Permanent {
		info, err :=
			os.Stat(folder.Path)

		if err == nil &&
			info.IsDir() {
			validPermanent =
				append(
					validPermanent,
					folder,
				)
		}
	}

	config.Permanent =
		validPermanent

	validRecent :=
		make([]Folder, 0)

	for _, folder := range config.Recent {
		info, err :=
			os.Stat(folder.Path)

		if err == nil &&
			info.IsDir() {
			validRecent =
				append(
					validRecent,
					folder,
				)
		}
	}

	config.Recent =
		validRecent
}

func saveConfig() {
	if configPath == "" {
		return
	}

	data, err :=
		json.MarshalIndent(
			config,
			"",
			"    ",
		)

	if err != nil {
		return
	}

	_ = os.WriteFile(
		configPath,
		data,
		0644,
	)
}

func displayName(
	path string,
) string {
	path =
		filepath.Clean(path)

	base :=
		filepath.Base(path)

	if base == "." ||
		base ==
			string(filepath.Separator) ||
		base == "" {
		return path
	}

	return base
}

func formatSize(
	size int64,
) string {
	units := []string{
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
	}

	value :=
		float64(size)

	for _, unit := range units {
		if value < 1024 {
			return fmt.Sprintf(
				"%.1f %s",
				value,
				unit,
			)
		}

		value /= 1024
	}

	return fmt.Sprintf(
		"%.1f PB",
		value,
	)
}
