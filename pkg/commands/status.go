package commands

import (
	"path/filepath"

	gogit "github.com/jesseduffield/go-git/v5"
)

type RebasingMode int

const (
	REBASE_MODE_NONE RebasingMode = iota
	REBASE_MODE_NON_INTERACTIVE
	REBASE_MODE_INTERACTIVE
)

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *Git) RebaseMode() RebasingMode {
	if c.gitDirFileExists("rebase-apply") {
		return REBASE_MODE_NON_INTERACTIVE
	}

	if c.gitDirFileExists("rebase-merge") {
		return REBASE_MODE_INTERACTIVE
	}

	return REBASE_MODE_NONE
}

func (c *Git) IsRebasing() bool {
	switch c.RebaseMode() {
	case REBASE_MODE_NON_INTERACTIVE, REBASE_MODE_INTERACTIVE:
		return true
	default:
		return false
	}
}

// IsMerging states whether we are still mid-merge
func (c *Git) IsMerging() bool {
	return c.gitDirFileExists("MERGE_HEAD")
}

func (c *Git) InNormalWorkingTreeState() bool {
	return !c.IsRebasing() && !c.IsMerging()
}

func (c *Git) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := c.repo.Worktree()
	return err == gogit.ErrIsBareRepository
}

func (c *Git) gitDirFileExists(path string) bool {
	result, err := c.GetOS().FileExists(filepath.Join(c.dotGitDir, path))
	if err != nil {
		// swallowing error
		c.log.Error(err)
	}

	return result
}
