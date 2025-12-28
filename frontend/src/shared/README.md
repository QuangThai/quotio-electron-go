# 📦 Shared Code

Thư mục này chứa code dùng chung giữa các process: **main**, **preload**, và **renderer**.

## 📁 Cấu trúc

```
shared/
├── README.md       # File này
├── constants.ts    # Hằng số chung của ứng dụng
└── utils.ts        # Các hàm tiện ích dùng chung
```

## 📄 Files

### `constants.ts`

Chứa các hằng số chung:

| Constant       | Mô tả                                        |
| -------------- | -------------------------------------------- |
| `APP_NAME`     | Tên ứng dụng: `"Quotio"`                     |
| `APP_VERSION`  | Phiên bản ứng dụng                           |
| `API_BASE_URL` | URL API backend: `http://localhost:8080/api` |
| `ROUTES`       | Object định nghĩa các route trong app        |

**Sử dụng:**

```typescript
import { APP_NAME, ROUTES, API_BASE_URL } from "../../shared/constants";

// Ví dụ
console.log(APP_NAME); // "Quotio"
console.log(ROUTES.DASHBOARD); // "/"
```

### `utils.ts`

Chứa các hàm tiện ích có thể dùng ở bất kỳ process nào:

| Function            | Mô tả                         |
| ------------------- | ----------------------------- |
| `formatDate(date)`  | Format ngày tháng theo locale |
| `formatNumber(num)` | Format số với dấu phân cách   |
| `sleep(ms)`         | Promise-based delay           |

**Sử dụng:**

```typescript
import { formatDate, formatNumber, sleep } from "../../shared/utils";

// Ví dụ
formatDate(new Date()); // "28/12/2024"
formatNumber(1234567); // "1,234,567"
await sleep(1000); // Đợi 1 giây
```

## 🔧 Nguyên tắc sử dụng

1. **Không import từ renderer process** - Code trong thư mục này phải hoạt động ở mọi process
2. **Không sử dụng DOM APIs** - Chỉ dùng JavaScript/TypeScript thuần
3. **Không sử dụng Node.js APIs trực tiếp** - Để đảm bảo tương thích với renderer process
4. **Export rõ ràng** - Mỗi function/constant cần được export riêng

## ➕ Thêm code mới

Khi thêm code mới vào thư mục này, hãy đảm bảo:

1. Code hoạt động ở cả main và renderer process
2. Có JSDoc comment mô tả chức năng
3. Cập nhật README này nếu cần thiết

**Ví dụ template:**

```typescript
/**
 * Mô tả chức năng của function
 * @param param1 - Mô tả tham số 1
 * @returns Mô tả giá trị trả về
 */
export const myNewFunction = (param1: string): string => {
  // Implementation
  return param1;
};
```
