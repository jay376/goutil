#include <vector>
#include <queue>
#include <stack>
using namespace std;

struct node {
  node* left;
  node* right;
  int value;
};

vector<int> travel(node* root) {
  vector<int> res;
  stack<node*> tmp;
  if (root == nullptr) {
    return res;
  }
  queue<node*> qu;
  qu.push(root);
  bool is_left = true;
  while (!qu.empty()) {
    int size = qu.size();
    while (size-->0) {
      auto n = qu.front();
      qu.pop();
      if (n->left) {
        qu.push(n->left);
      }
      if (n->right) {
        qu.push(n->right);
      }
      if (is_left) {
        res.push_back(n->value);
      } else {
        tmp.push(n);
      }
    }
    if (!is_left) {
      while (!tmp.empty()) {
        auto n = tmp.top();
        tmp.pop();
        res.push_back(n->value);
      }
      is_left = !is_left;
    }
  }

  return res;
}
