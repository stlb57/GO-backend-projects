// #include <bits/stdc++.h>
// using namespace std;


// int fun(vector<int> &arr, int j,int idx){
//     //1,7,2,23,4,5,6

//     //take - arr[idx] should be greater than j
//     if(idx==nums.size())return 0;
//     int take=0;
//     int not_take=0;
//     if(arr[idx]>j){
//         take=1+fun(arr,arr[idx],idx+1);
//         not_take=fun(arr,j,idx+1);
//     }
//     else{
//         not_take=fun(arr,j,idx+1);
//     }
//     return max(take,not_take);
// }

// int main() {
//      ios::sync_with_stdio(false);
//      cin.tie(nullptr);
//      int t;
//      cin>>t;
//      while(t--) {

//      }
// }


// // 2,4,5,100




// // int fun(vector<int> &arr, int j,int idx){
// //     //even sum
// //     // take -> odd + fun(odd sum)
// //     // take -> even + even (sum)

// //     // last step??
// //     // 1 no confirm -> odd/even 

// //     //j->0 

// //     int take,not_take;
// //     if(idx==arr.size()){
// //         return 0;
// //     }
// //     //if current element is even
// //     further_sum= fun(arr,idx+1);
// //     if(arr[idx]%2==0){
// //         if(further_sum%2==0){
// //             take=arr[idx]+further_sum;
// //         }
// //         else{
// //             not_take=further_sum;
// //         }
// //     }
// //     else{
// //         if(further_sum%2!=0){
// //             take=arr[idx]+further_sum;
// //         }
// //         else{
// //             not_take=further_sum;
// //         }
// //     }
// //     return max(take,not_take);

// // }

// // int main() {
// //      ios::sync_with_stdio(false);
// //      cin.tie(nullptr);
// //      int t;
// //      cin>>t;
// //      while(t--) {

// //      }
// // }


// // // 2,4,5,100
