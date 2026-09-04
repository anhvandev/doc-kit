package changelog

import "time"

// Add chèn mục mới lên đầu nhóm ngày tương ứng (tạo nhóm ở đầu nếu chưa có).
// Mục đầu nhóm cùng ngày, cùng Path, cùng Source, cách dưới window thì bị thay
// bằng mục mới: giờ mới, số dòng tính lại; hai tóm tắt thật khác nhau nối bằng
// "; " để không mất tóm tắt cũ (tóm tắt mới rỗng thì giữ tóm tắt cũ). Mục giữ chỗ (tóm tắt NoSummary, do hook agent ghi) được
// thay bất kể Source để tóm tắt thật của skill không nhân đôi dòng.
// Trả về mục thực sự được ghi.
func (f *File) Add(e Entry, window time.Duration) Entry {
	if e.Summary == "" {
		e.Summary = NoSummary
	}
	date := e.Time.Format(dateLayout)
	idx := -1
	for i, d := range f.Days {
		if d.Date == date {
			idx = i
			break
		}
	}
	if idx < 0 {
		f.Days = append([]Day{{Date: date}}, f.Days...)
		idx = 0
	}
	d := &f.Days[idx]
	if len(d.Entries) > 0 {
		last := d.Entries[0]
		delta := e.Time.Sub(last.Time)
		sameSource := last.Source == e.Source || last.Summary == NoSummary
		if last.Path == e.Path && sameSource && delta >= 0 && delta < window {
			switch {
			case e.Summary == NoSummary && last.Summary != NoSummary:
				e.Summary = last.Summary
			case last.Summary != NoSummary && last.Summary != e.Summary:
				e.Summary = last.Summary + "; " + e.Summary
			}
			d.Entries[0] = e
			return e
		}
	}
	d.Entries = append([]Entry{e}, d.Entries...)
	return e
}
