package hunter

const ()

func (h *Hunter) Dca(job DcaJob) {
	h.l.Debugw("Dca job", "job", job)
}
