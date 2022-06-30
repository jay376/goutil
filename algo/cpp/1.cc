#include "1.h"

int getRuntimeEndian() {
  int num = 1;
  int* pnum = &num;
  char* c = reinterpret_cast<char*>(pnum);
  if (*c == 1) {
    return 0;
  }
  return 1;
}
