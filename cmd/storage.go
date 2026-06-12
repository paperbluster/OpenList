/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/OpenListTeam/OpenList/v4/internal/bootstrap"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/spf13/cobra"
)

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage storage",
}

var disableStorageCmd = &cobra.Command{
	Use:   "disable [mount path]",
	Short: "Disable a storage by mount path",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("mount path is required")
		}
		mountPath := args[0]
		bootstrap.Init()
		defer bootstrap.Release()
		storage, err := db.GetStorageByMountPath(mountPath)
		if err != nil {
			return fmt.Errorf("failed to query storage: %+v", err)
		}
		storage.Disabled = true
		err = db.UpdateStorage(storage)
		if err != nil {
			return fmt.Errorf("failed to update storage: %+v", err)
		}
		utils.Log.Infof("Storage with mount path [%s] has been disabled from CLI", mountPath)
		fmt.Printf("Storage with mount path [%s] has been disabled\n", mountPath)
		return nil
	},
}

var deleteStorageCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a storage by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("id is required")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be a number")
		}

		if force, _ := cmd.Flags().GetBool("force"); force {
			fmt.Printf("Are you sure you want to delete storage with id [%d]? [y/N]: ", id)
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Delete operation cancelled.")
				return nil
			}
		}

		bootstrap.Init()
		defer bootstrap.Release()
		err = db.DeleteStorageById(uint(id))
		if err != nil {
			return fmt.Errorf("failed to delete storage by id: %+v", err)
		}
		utils.Log.Infof("Storage with id [%d] have been deleted from CLI", id)
		fmt.Printf("Storage with id [%d] have been deleted\n", id)
		return nil
	},
}

var listStorageCmd = &cobra.Command{
	Use:   "list",
	Short: "List all storages",
	RunE: func(cmd *cobra.Command, args []string) error {
		bootstrap.Init()
		defer bootstrap.Release()
		storages, _, err := db.GetStorages(1, -1)
		if err != nil {
			return fmt.Errorf("failed to query storages: %+v", err)
		}
		fmt.Printf("Found %d storages\n", len(storages))
		fmt.Printf("%-4s %-16s %-30s %-7s\n", "ID", "Driver", "Mount Path", "Enabled")
		for i := range storages {
			s := storages[i]
			enabled := "true"
			if s.Disabled {
				enabled = "false"
			}
			fmt.Printf("%-4d %-16s %-30s %-7s\n", s.ID, s.Driver, s.MountPath, enabled)
		}
		return nil
	},
}

func init() {

	RootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(disableStorageCmd)
	storageCmd.AddCommand(listStorageCmd)
	storageCmd.AddCommand(deleteStorageCmd)
	deleteStorageCmd.Flags().BoolP("force", "f", false, "Force delete without confirmation")
}
