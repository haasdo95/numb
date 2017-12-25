package versioning

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/haasdo95/numb/utils"

	git "github.com/libgit2/git2go"
	// git "gopkg.in/libgit2/git2go.v26"
)

func makeNumbSignature() *git.Signature {
	return &git.Signature{
		Name:  "numb",
		Email: "numb-team@gmail.com",
		When:  time.Now(),
	}
}

// CreateBranch creates a branch named "branchName"
func CreateBranch(repo *git.Repository, branchName string) {
	// retrieve most recent commit of master
	master, err := repo.LookupBranch("master", git.BranchLocal)
	utils.Check(err)
	masterHeadID := master.Target()
	masterHead, err := repo.LookupCommit(masterHeadID)
	utils.Check(err)

	// create new branch
	_, err = repo.CreateBranch(branchName, masterHead, false)
	if err != nil {
		fmt.Println("Failed to create the branch reserved for 'numb'")
		fmt.Println("Probably you have a branch named 'numb' already")
		fmt.Println("Run 'git branch' to figure out")
	}
}

func makeCommitToNumb(commitMessage string) *git.Oid {
	repo, err := git.OpenRepository(".git")
	utils.Check(err)
	idx, err := repo.Index()
	utils.Check(err)
	treeID, err := idx.WriteTree()
	utils.Check(err)
	tree, err := repo.LookupTree(treeID)
	utils.Check(err)

	head, _ := repo.Head()
	headID := head.Target()
	headCommit, err := repo.LookupCommit(headID)
	utils.Check(err)

	// commit to head
	numbSignature := makeNumbSignature()
	defaultSignature, err := repo.DefaultSignature()
	utils.Check(err)
	commitID, err := repo.CreateCommit("refs/heads/numb", defaultSignature, numbSignature, commitMessage, tree, headCommit)
	utils.Check(err)
	return commitID
}

// FlashCommit does the following:
// 1. stash uncommited changes on current branch. remember current branch.
// 2. checkout numb branch, APPLY the stashed
// 3. resolve conflicts by always accepting incoming changes; add -A & commit
func FlashCommit(params string) (*git.Oid, error) {
	repo, err := git.OpenRepository(".git")
	utils.Check(err)
	defer repo.Free()
	// 1
	// remember old HEAD
	head, err := repo.Head()
	if !head.IsBranch() {
		fmt.Println("You may be in 'HEAD Detached' state")
		return nil, errors.New("Bad Head")
	}
	oldBranch := head.Branch()
	oldBranchName, err := oldBranch.Name()
	utils.Check(err)

	// stash stuff
	stasher := makeNumbSignature()
	stashID, err := repo.Stashes.Save(stasher, "Stashing the Uncommitted on Working Branch", git.StashIncludeUntracked)
	if stashID == nil || err != nil {
		println("Nothing Stashed!")
		checkoutCmd := exec.Command("git", "checkout", "numb")
		checkoutCmd.Run()
		mergeCmd := exec.Command("git", "merge", "-X", "theirs", oldBranchName)
		mergeCmd.Run()

		addAllCmd := exec.Command("git", "add", "-A")
		addAllCmd.Run()
		commit := makeCommitToNumb("Trained at " + time.Now().String())
		checkoutBackCmd := exec.Command("git", "checkout", oldBranchName)
		checkoutBackCmd.Run()
		return commit, nil
	}

	// make sure to checkout back
	rewind := func() {
		checkoutBackCmd := exec.Command("git", "checkout", oldBranchName)
		checkoutBackCmd.Run()

		stashPopCmd := exec.Command("git", "stash", "pop")
		stashPopCmd.Run()
	}
	defer rewind()

	// 2
	// checkout to numb
	checkoutCmd := exec.Command("git", "checkout", "numb")
	checkoutCmd.Run()
	// apply the stashed
	stashAppCmd := exec.Command("git", "stash", "apply")
	stashAppCmd.Run()

	// 3
	// resolve conflicts TODO:
	println("Trying to resovle conflicts")
	resolveCmd := exec.Command("bash", "-c", "grep -lr '<<<<<<<' . | xargs git checkout --theirs")
	resolveCmd.Run()

	addAllCmd := exec.Command("git", "add", "-A")
	addAllCmd.Run()
	commitMessage := "Trained at " + time.Now().String() + "\n"
	commitMessage += "With Params: " + params
	commit := makeCommitToNumb(commitMessage)

	println("Commit Made: ", commit.String())

	return commit, nil
}
