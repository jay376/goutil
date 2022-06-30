#include <iostream>
#include <ostream>
#include "risk.h"

// All Rights Reserved.
// Author : zhangfangjie (f22jay@163.com)
// Date 2022/06/13 09:41:25
// Breif :

int main() {
  Risk::Policies policies;
  Risk r(policies);
  Order order;
  std::cout << r.check(order) << std::endl;
}
