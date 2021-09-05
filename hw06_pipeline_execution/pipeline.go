package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for i := 0; i < len(stages); i++ {
		out = intercept(done, stages[i](out))
	}

	return out
}

func intercept(done In, out Out) Out {
	outCopy := make(Bi, 1)

	go func() {
		defer close(outCopy)
		for {
			select {
			case <-done:
				return
			case value, ok := <-out:
				if ok {
					select {
					case outCopy <- value:
					case <-done:
						return
					}

					break
				}

				return
			}
		}
	}()

	return outCopy
}
