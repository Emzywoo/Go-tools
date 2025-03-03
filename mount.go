package main

import (
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

const (
    mediaPath = "/media/emmadev"
    diskByLabelPath = "/dev/disk/by-label"
)

func main() {
    log.Println("Starting the mount process")

    // Ensure the base media path exists
    err := os.MkdirAll(mediaPath, 0755)
    if err != nil {
        log.Fatalf("Failed to create base media path: %v", err)
    }
    log.Printf("Base media path ensured: %s", mediaPath)

    // Read the disk by label directory
    labels, err := ioutil.ReadDir(diskByLabelPath)
    if err != nil {
        log.Fatalf("Failed to read %s: %v", diskByLabelPath, err)
    }

    for _, label := range labels {
        decodedLabel := decodeLabel(label.Name())
        if isExternalDrive(decodedLabel) {
            mountDrive(label.Name(), decodedLabel)
        } else {
            log.Printf("Skipping non-external drive: %s", decodedLabel)
        }
    }

    log.Println("Mount process completed")
}

func decodeLabel(label string) string {
    return strings.ReplaceAll(label, "\\x20", " ")
}

func isExternalDrive(label string) bool {
    ignoredLabels := []string{"System Reserved", "Windows", "Recovery"}
    for _, ignored := range ignoredLabels {
        if strings.EqualFold(label, ignored) {
            return false
        }
    }
    return true
}

func mountDrive(encodedLabel, decodedLabel string) {
    devicePath := filepath.Join(diskByLabelPath, encodedLabel)
    mountPoint := filepath.Join(mediaPath, decodedLabel)

    log.Printf("Attempting to mount drive with label: %s", decodedLabel)

    // Create mount point directory if it doesn't exist
    if err := os.MkdirAll(mountPoint, 0755); err != nil {
        log.Printf("Failed to create mount point directory for %s: %v", decodedLabel, err)
        return
    }

    // Attempt to mount the device
    cmd := exec.Command("mount", devicePath, mountPoint)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Failed to mount %s: %v\nCommand output: %s", decodedLabel, err, string(output))
    } else {
        log.Printf("Successfully mounted %s to %s", decodedLabel, mountPoint)
    }
}
