# 🧩 Shared Components

Thư mục này chứa các **React components dùng chung** trên nhiều trang/màn hình trong ứng dụng.

## 📁 Cấu trúc

```
shared/
├── README.md           # File này
└── status-badges.tsx   # Components hiển thị status badges
```

## 📄 Components

### `status-badges.tsx`

Chứa các components để hiển thị trạng thái dưới dạng badges.

#### `renderAccountStatusBadge(status)`

Hiển thị trạng thái của account/provider.

| Status                    | Badge        | Variant           |
| ------------------------- | ------------ | ----------------- |
| `active` hoặc `undefined` | Active       | `success` (xanh)  |
| `rate_limited`            | Rate Limited | `warning` (vàng)  |
| `cooldown`                | Cooldown     | `danger` (đỏ)     |
| Khác                      | Disabled     | `secondary` (xám) |

**Sử dụng:**

```tsx
import { renderAccountStatusBadge } from "../shared/status-badges";

// Trong component
<div>{renderAccountStatusBadge(provider.status)}</div>;
```

#### `renderAgentStatusBadge(installed, configured, hasError)`

Hiển thị trạng thái của agent.

| Điều kiện           | Badge         | Variant           |
| ------------------- | ------------- | ----------------- |
| `hasError = true`   | Config Error  | `danger` (đỏ)     |
| `configured = true` | Configured    | `success` (xanh)  |
| `installed = true`  | Installed     | `warning` (vàng)  |
| Tất cả `false`      | Not Installed | `secondary` (xám) |

**Sử dụng:**

```tsx
import { renderAgentStatusBadge } from "../shared/status-badges";

// Trong component
<div>
  {renderAgentStatusBadge(agent.installed, agent.configured, agent.hasError)}
</div>;
```

## 🔧 Nguyên tắc

1. **Tái sử dụng** - Components ở đây phải được dùng ở ít nhất 2 nơi
2. **Đơn giản** - Mỗi component chỉ làm một việc
3. **Props rõ ràng** - Định nghĩa types cho tất cả props
4. **Import từ UI** - Sử dụng components từ `../ui/` thay vì tự implement

## ➕ Thêm component mới

Khi thêm component mới, hãy:

1. Đảm bảo component được dùng ở nhiều nơi
2. Thêm JSDoc comments
3. Cập nhật README này

**Ví dụ template:**

```tsx
import { SomeUIComponent } from "../ui/some-component";

/**
 * Mô tả component này làm gì
 * @param props.value - Mô tả prop
 */
export function MySharedComponent({ value }: { value: string }) {
  return <SomeUIComponent>{value}</SomeUIComponent>;
}
```

## 🔗 Liên quan

- **UI Components**: `../ui/` - Components cơ bản (Button, Badge, Card, etc.)
- **Shared Utils**: `../../../shared/` - Utilities dùng chung cho cả app
