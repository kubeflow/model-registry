package ptr

func In[E any](e *E) E {
	if e == nil {
		return *new(E)
	}

	return *e
}

func Of[E any](e E) *E {
	return &e
}
