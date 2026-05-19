# BACK-END SPECS (Derived from `docs/requirement.md`)

Tài liệu này đặc tả phần back-end hỗ trợ toàn bộ yêu cầu REQ-001 đến REQ-016.

## 0) Traceability — Back-end support cho toàn bộ requirement

### 0.1 Nhóm Front-end (REQ-001 đến REQ-008)

- **REQ-001 (Header & Điều hướng):** Cung cấp API tìm kiếm sản phẩm theo từ khóa; API cấu hình menu/topbar; API cấu hình logo.
- **REQ-002 (Banner & Danh mục nổi bật):** Cung cấp API CRUD banner; API danh mục nổi bật + icon; API liên kết icon -> danh mục.
- **REQ-003 (Khối sản phẩm theo danh mục):** API trả block sản phẩm theo danh mục, hỗ trợ giá gốc/giá bán/% giảm, nhãn sản phẩm, trạng thái còn hàng để thêm giỏ.
- **REQ-004 (Footer, Tin tức mới, Floating):** API cấu hình footer theo cột; API tin tức mới nhất có giới hạn số lượng; API cấu hình nút floating (bật/tắt, link).
- **REQ-005 (Trang Giới thiệu):** API quản lý trang giới thiệu, gallery, và danh sách tin liên quan.
- **REQ-006 (Danh sách sản phẩm):** API cây danh mục, lọc giá, lọc thương hiệu, sorting, phân trang; API dữ liệu sản phẩm dạng grid.
- **REQ-007 (Chi tiết sản phẩm):** API chi tiết sản phẩm (ảnh, thumbnail, giá, SKU, quà tặng); API tab nội dung và tổng hợp số lượng đánh giá.
- **REQ-008 (Tin tức & Liên hệ):** API danh sách bài viết dạng lưới; API thông tin liên hệ + map config; API FAQ; API submit form tư vấn.

### 0.2 Nhóm Tính năng Khách hàng (REQ-009 đến REQ-010)

- **REQ-009 (Mua hàng & Giỏ hàng):** API giỏ hàng (thêm/sửa số lượng/xóa/xem); API tạo đơn hàng/yêu cầu đặt hàng.
- **REQ-010 (Tương tác & Hỗ trợ):** API tiếp nhận form tư vấn từ trang chi tiết và liên hệ; API cấu hình các kênh hỗ trợ trực tuyến (live chat, gọi điện, Zalo, Facebook).

### 0.3 Nhóm Quản trị (REQ-011 đến REQ-014)

- **REQ-011:** Account & phân quyền.
- **REQ-012:** Quản lý bài viết/trang tĩnh + lịch xuất bản.
- **REQ-013:** Quản lý sản phẩm/danh mục/tags/thuộc tính.
- **REQ-014:** Media + Form builder + quản lý phản hồi.

### 0.4 Nhóm Phi chức năng & Kỹ thuật (REQ-015 đến REQ-016)

- **REQ-015:** Vận hành trong giới hạn hạ tầng, bắt buộc HTTPS/SSL, hỗ trợ dữ liệu SEO.
- **REQ-016:** Dữ liệu/API tương thích đa thiết bị; hỗ trợ cấu hình theme; hỗ trợ nhập liệu khởi tạo.

## 1) SPEC-BE-011 — Quản lý Tài khoản & Phân quyền

- **Nguồn yêu cầu:** REQ-011
- **Mục tiêu:** Quản lý người dùng trong admin và kiểm soát truy cập theo vai trò.

### 1.1 Phạm vi chức năng
- Thêm/Sửa/Xóa/Khóa tài khoản người dùng.
- Phân quyền tối thiểu 5 cấp độ: `Administrator`, `Editor`, `Author`, `Contributor`, `Subscriber`.
- Người dùng đăng nhập có thể cập nhật hồ sơ cá nhân.

### 1.2 Quy tắc nghiệp vụ
- Mỗi người dùng có đúng 1 vai trò chính tại một thời điểm.
- Tài khoản bị khóa không được đăng nhập.
- Chỉ `Administrator` được quyền thay đổi vai trò hoặc khóa/mở khóa tài khoản.
- Người dùng chỉ được sửa hồ sơ của chính mình, trừ `Administrator`.

### 1.3 Tiêu chí chấp nhận
- Tạo mới người dùng thành công khi dữ liệu hợp lệ.
- Cập nhật/xóa/khóa người dùng phản ánh đúng trên danh sách quản trị.
- Phân quyền đúng theo vai trò; thao tác vượt quyền bị từ chối.
- Người dùng đăng nhập cập nhật hồ sơ thành công và dữ liệu được lưu.

---

## 2) SPEC-BE-012 — Quản lý Nội dung (Tin bài & Trang tĩnh)

- **Nguồn yêu cầu:** REQ-012
- **Mục tiêu:** Quản trị nội dung blog và các trang tĩnh.

### 2.1 Phạm vi chức năng
- Quản lý bài viết: Thêm/Sửa/Xóa.
- Trường thông tin bài viết: tiêu đề, nội dung, ảnh đại diện, danh mục, thẻ (tags).
- Hỗ trợ lên lịch xuất bản.
- Quản lý trang tĩnh: Thêm/Sửa/Xóa (ví dụ: Giới thiệu, Liên hệ).
- Cho phép chọn template cho từng trang tĩnh.

### 2.2 Quy tắc nghiệp vụ
- Bài viết có trạng thái: `draft`, `scheduled`, `published`.
- Bài viết `scheduled` tự chuyển `published` khi đến thời điểm cấu hình.
- Không cho xuất bản nếu thiếu các trường bắt buộc (ít nhất: tiêu đề, nội dung chính).
- Trang tĩnh phải có duy nhất 1 template đang áp dụng.

### 2.3 Tiêu chí chấp nhận
- CRUD bài viết hoạt động đúng với dữ liệu hợp lệ.
- Lên lịch xuất bản thực thi đúng thời điểm.
- CRUD trang tĩnh hoạt động đúng.
- Gán/chuyển đổi template cho trang tĩnh thành công.

---

## 3) SPEC-BE-013 — Quản lý Sản phẩm & Danh mục

- **Nguồn yêu cầu:** REQ-013
- **Mục tiêu:** Quản trị dữ liệu sản phẩm và danh mục cho website bán hàng.

### 3.1 Phạm vi chức năng
- CRUD sản phẩm với trường tối thiểu: tiêu đề, mô tả, ảnh, giá, SKU.
- Quản lý cây danh mục sản phẩm.
- Gắn thẻ (tags) cho sản phẩm.
- Cấu hình thuộc tính riêng cho sản phẩm (ví dụ: thương hiệu, biến thể, thông số).

### 3.2 Quy tắc nghiệp vụ
- SKU là duy nhất theo từng sản phẩm.
- Giá bán phải lớn hơn hoặc bằng 0.
- Danh mục có thể phân cấp cha-con.
- Một sản phẩm có thể thuộc nhiều danh mục và nhiều tags.

### 3.3 Tiêu chí chấp nhận
- CRUD sản phẩm và danh mục hoạt động ổn định.
- Hệ thống chặn SKU trùng.
- Thuộc tính mở rộng được lưu và truy xuất đúng.
- Liên kết sản phẩm với danh mục/tags phản ánh đúng ở dữ liệu.

---

## 4) SPEC-BE-014 — Quản lý Kho lưu trữ & Tương tác

- **Nguồn yêu cầu:** REQ-014
- **Mục tiêu:** Quản lý media, form liên hệ và phản hồi khách hàng.

### 4.1 Phạm vi chức năng
- Media Library: tải lên, lưu trữ, xem danh sách, xóa tài nguyên (ảnh/video/audio).
- Form Builder: tạo và tùy chỉnh các trường form liên hệ.
- Quản lý phản hồi: xem danh sách form submit, cập nhật trạng thái xử lý.

### 4.2 Quy tắc nghiệp vụ
- Chỉ cho phép định dạng media nằm trong danh sách cho phép.
- Có kiểm tra kích thước tệp khi upload theo cấu hình hệ thống.
- Mỗi phản hồi form có trạng thái xử lý (ví dụ: `new`, `in_progress`, `done`).
- Dữ liệu submit phải gắn timestamp và thông tin form nguồn.

### 4.3 Tiêu chí chấp nhận
- Upload và xóa media thành công với tệp hợp lệ.
- Tạo/chỉnh sửa form và áp dụng cho trang liên hệ thành công.
- Dữ liệu phản hồi được lưu đầy đủ và cập nhật trạng thái được.

---

## 5) Phi chức năng áp dụng cho khối Back-end

Dựa trên REQ-015 và REQ-016, phần back-end cần đảm bảo:

- Hỗ trợ vận hành ổn định trên môi trường hosting/server có hạn mức lưu trữ phù hợp.
- Hỗ trợ HTTPS/SSL khi triển khai.
- Cung cấp dữ liệu phục vụ SEO on-page (meta, URL thân thiện, heading/alt liên quan nội dung).
- API/back-end tương thích tốt với giao diện responsive đa thiết bị.
