
// All Rights Reserved.
// Author : zhangfangjie (f22jay@163.com)
// Date 2021/02/02 14:59:06
// Breif :

#include <stack>
#include <vector>
#include <queue>
#include <unordered_map>
using namespace std;

struct Node {
    Node* left;
    Node* right;
    int value;
};

vector<int> inorder(Node* root) {
    vector<int> res;
    stack<Node*> sk;
    if (root == nullptr) {
        return res;
    }

    while (!sk.empty() || root != nullptr) {
        if (root != nullptr) {
            sk.push(root);
            root = root->left;
        } else {
            root = sk.top();
            sk.pop();
            res.push_back(root->value);
            root = root->right;
        }
    }
    return res;
}

vector<int> preOrder(Node* root) {
    vector<int> res;
    stack<Node*> sk;
    if (root == nullptr) {
        return res;
    }

    while (!sk.empty() || root != nullptr) {
        if (root != nullptr) {
            res.push_back(root->value);
            sk.push(root);
            root = root->left;
        } else {
            root = sk.top();
            sk.pop();
            root = root->right;
        }
    }
    return res;
}

// ListNode* mergeKLists(vector<ListNode*>& lists) {
//     priority_queue<ListNode*,dequeue<ListNode*>> heap;
//     for (auto&& it : lists) {
//         heap.push(lists);
//     }
//     ListNode* head = new ListNode(0);
//     ListNode** p = &head.next;
//     while (!heap.empty()) {
//         auto node = heap.top();
//         heap.pop();
//         *p = node;
//         p = &(node->next);
//         node = node->next;
//         if (node != nullptr) {
//             heap.push(node);
//         }
//     }
//     *p =nullptr
//     return head->next;
// }

vector<int> postOrder(Node* root) {
    vector<int> res;
    stack<Node*> sk;
    if (root == nullptr) {
        return res;
    }
    Node  *node = root, *prev = nullptr;

    while (node != nullptr || !sk.empty()) {
        if (node != nullptr) {
            sk.push(node);
            node = node->left;
        } else {
            node = sk.top();
            if (node->right != nullptr && node->right != prev) {
                node = node->right;
            } else {
                res.push_back(node->value);
                prev = node;
                sk.pop();
            }
        }
    }
}

class Solution {
public:
    int firstMissingPositive(vector<int>& nums) {
        int size = nums.size();
        for (int i = 0; i < size;) {
            int dst = nums[i];
            if (dst > size || dst <= 0 || dst == (i+1) ) {
                i++;
                continue;
            }
            int tmp = nums[dst-1];
            nums[dst-1] = dst;
            nums[i] = tmp;
        }
        int lost = 0;
        for (int i = 0; i < size; i++) {
            if (nums[i] != (i+1)) {
                lost = i+1;
                break;
            }
        }

        return lost;
    }
};

struct ListNode {
  int val;
  ListNode* next;
};

struct compare {
  bool operator() (ListNode* left, ListNode* right) {
    return left->val < right->val;
  }
};

ListNode* mergeKLists(vector<ListNode*>& lists) {
  std::priority_queue<ListNode*, vector<ListNode*>, compare> heap(lists.begin(), lists.end());
  ListNode* head = new ListNode();
  ListNode** next  = &head->next;
  while (!heap.empty()) {
    auto node = heap.top();
    *next = node;
    next  = &node->next;
    heap.pop();
    if (node->next != nullptr) {
      heap.push(node->next);
    }
  }

  auto ret = head->next;
  delete head;
  return ret;
}


struct List {
  List* prev;
  List* next;
  int key;
};
void add(List* head, List* node) {
  node->next = head->next;
  node->prev = head;
  head->next->prev = node;
  head->next = node;
}

void del(List* node) {
  node->next->prev = node->prev;
  node->prev->next = node->next;
}

struct entry {
  int key;
  int value;
  List* node;
};

class LRUCache {
 private:
  int cap_;
  int length_;
  std::unordered_map<int, entry*> hmap_;
  List* head_;
 public:

    LRUCache(int capacity):cap_(capacity), length_(0) {
      head_ = new List();
      head_->next = head_;
      head_->prev = head_;
    }

    int get(int key) {
      auto it = hmap_.find(key);
      if (it == hmap_.end()) {
        return -1;
      }
      del(it->second->node);
      add(head_, it->second->node);
      return it->second->value;
    }

    void put(int key, int value) {
      entry* e = nullptr;
      auto  it = hmap_.find(key);
      if (it == hmap_.end()) {
        if (length_ == cap_) {
          auto node = head_->prev;
          del(node);
          hmap_.erase(node->key);
          delete(node);
          length_--;
        }
        auto node = new List();
        node->next  = node;
        node->prev  = node;
        node->key = key;
        e = new entry{key, value, node};
        length_++;
      } else {
        e->value = value;
      }
      del(e->node);
      add(head_, e->node);
      hmap_[key] = e;
    }
};

// https://leetcode.cn/problems/median-of-two-sorted-arrays/submissions/
// 中位数
class Solution {
 public:
  double findK(vector<int>& nums1, vector<int>& nums2, int k) {
    int l1 = 0, l2 = 0;
    int r1 =  nums1.size() - 1, r2 = nums2.size() - 1;
    while (true) {
      if (l1 > r1) {
        return nums2[l2+k-1];
      }
      if (l2 > r2) {
        return nums1[l1+k-1];
      }
      if (k==1) {
        return min(nums1[l1], nums2[l2]);
      }
      int m1 = min(l1 + k/2 -1, r1);
      int m2 = min(l2 + k/2 -1, r2);
      if (nums1[m1] <= nums2[m2]) {
        k -= (m1-l1+1);
        l1 = m1+1;
      }
      if (nums1[m1] > nums2[m2]) {
        k -= (m2-l2+1);
        l2 = m2+1;
      }
      if (k<=0) {
        return max(nums1[m1], nums2[m2]);
      }
    }
    return 0;
  }
    double findMedianSortedArrays(vector<int>& nums1, vector<int>& nums2) {
      int size = nums1.size() + nums2.size();
      if (size % 2 == 0) {
        return (findK(nums1, nums2, size/2)  + findK(nums1, nums2, size/2 + 1)) / 2;
      }
      return findK(nums1, nums2, size/2+1);
    }
};
