#!/bin/sh
# Do dk cài: chặn commit khi docs/ có file đổi mà chưa có dòng changelog.
# Thiếu dk thì cho qua có chủ đích; xem README mục "Cài vào dự án".
command -v dk >/dev/null 2>&1 || { echo "dk chưa cài, bỏ qua kiểm tra changelog"; exit 0; }
exec dk --cwd "__DK_CWD__" changelog pending
