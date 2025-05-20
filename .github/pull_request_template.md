name: Pull Request
description: Tạo pull request để đóng góp cho dự án
title: "[PR]: "
body:
  - type: markdown
    attributes:
      value: |
        Cảm ơn bạn đã tạo pull request! Vui lòng đảm bảo bạn đã điền đầy đủ thông tin dưới đây.
  - type: dropdown
    id: type
    attributes:
      label: Loại thay đổi
      description: Loại thay đổi mà PR này bao gồm
      options:
        - feat (tính năng mới)
        - fix (sửa lỗi)
        - docs (tài liệu)
        - style (định dạng, thiếu dấu chấm phẩy, v.v.)
        - refactor (thay đổi code không thêm tính năng hoặc sửa lỗi)
        - perf (cải thiện hiệu suất)
        - test (thêm hoặc sửa test)
        - chore (thay đổi công cụ build, v.v.)
    validations:
      required: true
  - type: textarea
    id: description
    attributes:
      label: Mô tả thay đổi
      description: Mô tả chi tiết về các thay đổi trong PR này
      placeholder: Mô tả chi tiết về những gì bạn đã thay đổi và lý do
    validations:
      required: true
  - type: textarea
    id: related-issues
    attributes:
      label: Issues liên quan
      description: Liên kết đến issues mà PR này giải quyết
      placeholder: "Fixes #123, Resolves #456"
  - type: textarea
    id: testing
    attributes:
      label: Testing thực hiện
      description: Mô tả các tests bạn đã thực hiện để xác minh thay đổi của mình
      placeholder: Unit tests, integration tests, v.v.
    validations:
      required: true
  - type: checkboxes
    id: checklist
    attributes:
      label: Checklist
      options:
        - label: Tôi đã thêm tests cho code mới
          required: true
        - label: Tôi đã cập nhật tài liệu nếu cần thiết
          required: true
        - label: Code của tôi tuân thủ tiêu chuẩn code của dự án
          required: true
        - label: Tôi đã chạy `go fmt` và `go vet`
          required: true
