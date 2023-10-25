package git

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Merge attempts to merge ref onto HEAD. Currently only supports fast-forward merges
func (r *Repository) Merge(ref plumbing.Reference, opts MergeOptions) error {
	switch opts.Mode {
	case MergeAbort:
		return errors.New("unsupported merge mode: MergeAbort")
	case MergeContinue:
		return errors.New("unsupported merge mode: MergeContinue")
	case MergeDefault:
		return r.merge(ref, opts)
	default:
		return fmt.Errorf("unknown merge mode: %d", opts.Mode)
	}
}

func (r *Repository) merge(ref plumbing.Reference, opts MergeOptions) error {
	head, err := r.Head()
	if err != nil {
		return err
	}

	ff, err := isFastForward(r.Storer, head.Hash(), ref.Hash())
	if err != nil && err != storer.ErrStop {
		return fmt.Errorf("fast forward is not possible: %w", err)
	}

	if !ff {
		return errors.New("non fast-forward merges are not supported yet")
	}

	switch opts.FastForward {
	case FastForward:
		return r.Storer.SetReference(
			plumbing.NewHashReference(head.Name(), ref.Hash()))
	}

	return fmt.Errorf("unsupported fast-forward option: %s", opts.FastForward)
}
