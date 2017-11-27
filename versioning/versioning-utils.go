package versioning

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/user/numb/utils"

	"github.com/libgit2/git2go"
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

func addAllAndCommitOnNumb(repo *git.Repository, commitMessage string) (*git.Oid, bool, error) {
	// get HEAD
	head, err := repo.Head()
	// check head
	if !head.IsBranch() {
		fmt.Println("You may be in 'HEAD Detached' state")
		return nil, false, errors.New("Bad Head")
	} else {
		headBranch := head.Branch()
		headName, err := headBranch.Name()
		utils.Check(err)
		if headName != "numb" { // not on numb
			fmt.Println("You are currently on branch: ", headName)
			fmt.Println("This function needs you to be on numb branch")
			return nil, false, errors.New("Wrong Current Branch")
		}
	}

	// git add -A
	idx, err := repo.Index()
	utils.Check(err)
	hasConflict := idx.HasConflicts()
	err = idx.AddAll([]string{}, git.IndexAddDefault, nil)
	utils.Check(err)
	treeID, err := idx.WriteTree()
	utils.Check(err)
	tree, err := repo.LookupTree(treeID)
	utils.Check(err)
	err = idx.Write()
	utils.Check(err)

	utils.Check(err)
	headID := head.Target()
	headCommit, err := repo.LookupCommit(headID)
	utils.Check(err)

	// commit to head
	numbSignature := makeNumbSignature()
	defaultSignature, err := repo.DefaultSignature()
	utils.Check(err)
	commitID, err := repo.CreateCommit("refs/heads/numb", defaultSignature, numbSignature, commitMessage, tree, headCommit)
	utils.Check(err)
	return commitID, hasConflict, nil
}

func checkoutBranch(repo *git.Repository, branchName string) error {
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe,
	}
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if branch == nil || err != nil {
		return errors.New("Failed to lookup branch: " + branchName)
	}
	commit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		log.Println("Failed to lookup commit")
		return err
	}
	tree, err := repo.LookupTree(commit.TreeId())
	if err != nil {
		log.Println("Failed to lookup tree")
		return err
	}
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		log.Println("Failed to checkout tree")
		return err
	}
	repo.SetHead("refs/heads/" + branchName)
	return nil
}

// FlashCommit does the following:
// 1. stash uncommited changes on current branch. remember current branch.
// 2. checkout numb branch, APPLY the stashed
// 3. resolve conflicts by always accepting incoming changes; add -A & commit
func FlashCommit(repo *git.Repository) (*git.Oid, error) {
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
	repo.Stashes.Save(stasher, "Stashing the Uncommitted on Working Branch", git.StashIncludeUntracked)

	// 2
	// checkout to numb
	err = checkoutBranch(repo, "numb")
	utils.Check(err)
	// apply the stashed
	defaultStashOptions, err := git.DefaultStashApplyOptions()
	utils.Check(err)
	err = repo.Stashes.Apply(0, defaultStashOptions)
	utils.Check(err)

	// 3
	// resolve conflicts
	// TODO: Haven't figured out the safe way to resolve conflicts.
	// For now I'll just yell at the user to manually fix them
	// add -A & commit
	commit, hasConflict, err := addAllAndCommitOnNumb(repo, "Trained at "+time.Now().String())
	utils.Check(err)

	if hasConflict {
		fmt.Println("Conflict(s) occurred on numb branch")
		fmt.Println("checkout to numb branch and fix conflict before doing anything else")
		fmt.Println("Most times, simply keep hitting 'accept incoming changes' is a safe bet")
	}

	// flash back
	err = checkoutBranch(repo, oldBranchName)
	utils.Check(err)
	err = repo.Stashes.Pop(0, defaultStashOptions)
	return commit, nil
}
