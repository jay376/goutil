// All Rights Reserved.
// Author : zhangfangjie (f22jay@163.com)
// Date 2022/06/13 09:42:36
// Breif :
#ifndef RISK_H_
#define RISK_H_

#include <map>
#include <stdint.h>
#include <mutex>
#include <memory>

struct Order {
  char instrumentId[20];  //股票代码，如: 111111.SH, 222222.SH
  uint64_t accountId; //报单账户（对应某个产品）
  double price;       //报单价格
  int32_t qty;        //报单数量
  // 其他字段...
};

class ProductPolicy {
 public:
  virtual bool check(const Order& order) = 0;
};

class AProductPolicy: public ProductPolicy {
  bool check(const Order& order) {
    // xxxx
    return true;
  }
};

class BProductPolicy: public ProductPolicy {
  bool check(const Order& order) {
    // xxxx
    return true;
  }
};

class CProductPolicy: public ProductPolicy {
  bool check(const Order& order) {
    // xxxx
    return true;
  }
};

class Risk {
 public:
  typedef std::map<uint64_t, std::shared_ptr<ProductPolicy>> Policies;
  Risk(const Policies& policies): policies_(policies) {}
  // preCheck 做些公共检查
  bool preCheck(const Order& order) {
    // 所有产品（A，B，C）都不能交易代码为123456.SH
    return true;
  }

  bool check(const Order& order) {
    // 1. 做公共check
    if (!preCheck(order)) {
      return false;
    }

    // 2. 做具体产品check
    mu_.lock();
    auto it = policies_.find(order.accountId);
    mu_.unlock();
    if (it != policies_.end()) {
      return  it->second->check(order);
    }
    return false;
  }

  bool addPolicy(uint64_t type, std::shared_ptr<ProductPolicy> policy) {
    std::lock_guard<std::mutex> lc(mu_);
    policies_[type] = policy;
  }

 private:
  Policies policies_;
  std::mutex mu_;
};

#endif
