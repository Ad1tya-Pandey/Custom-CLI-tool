package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

func GetCPUUsage() string {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return fmt.Sprintf("Error getting CPU Usage: %s", err)
	}
	return fmt.Sprintf("CPU Usage: %.2f%%", percentages[0])
}

///////////////////////////////////////

func GetMemoryUsage() string {
	// Fetch memory usage statistics
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "Error getting Memory Usage: " + err.Error()
	}

	// Create a new table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Metric", "Value"})

	// Add memory usage metrics to the table
	table.Append([]string{"Total", formatBytes(vmStat.Total)})
	table.Append([]string{"Free", formatBytes(vmStat.Free)})
	table.Append([]string{"Used", formatBytes(vmStat.Used)})
	table.Append([]string{"Used Percent", formatPercent(vmStat.UsedPercent)})

	// Render the table
	var buf bytes.Buffer
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.Render()

	// Write the table's output to the buffer
	tableWriter := tablewriter.NewWriter(&buf)
	tableWriter.SetBorder(false)
	tableWriter.SetColumnSeparator("")
	tableWriter.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	tableWriter.Render()

	// Return the table's output as a string
	return buf.String()
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div := uint64(unit)
	exp := 0
	for {
		if bytes < div*unit {
			return fmt.Sprintf("%.2f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
		}
		div *= unit
		exp++
	}
}

func formatPercent(percent float64) string {
	return fmt.Sprintf("%.2f%%", percent)
}

////////////////////////////

func GetDiskUsage() string {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return fmt.Sprintf("Error getting Disk Usage: %s", err)
	}
	return fmt.Sprintf("Disk Usage: Total: %v, Free: %v, UsedPercent: %.2f%%", diskStat.Total, diskStat.Free, diskStat.UsedPercent)
}

/////////////////////////////

func DisplayProcessInfo() {
	// Create a tabwriter to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	// Print the header
	fmt.Fprintln(w, "PID\tName\tCPU%\tMemory%")

	// Fetch a list of running processes
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
		return
	}

	// Sort processes by PID
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Pid < processes[j].Pid
	})

	// Print information about each process
	for _, p := range processes {
		printProcessInfo(w, p)
	}
}

//////////////////////////////////

func printProcessInfo(w *tabwriter.Writer, p *process.Process) {
	name, _ := p.Name()
	cpuPercent, _ := p.CPUPercent()
	memoryPercent, _ := p.MemoryPercent()

	fmt.Fprintf(w, "%d\t%s\t%.2f\t%.2f\n", p.Pid, name, cpuPercent, memoryPercent)
}

//////////////////////////////////

func DisplayNetworkUsage() {
	// Get network usage statistics
	netStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("Error getting network usage:", err)
		return
	}

	// Create a tabwriter to format the output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	// Print the header
	fmt.Fprintln(w, "Interface\tRX Bytes\tTX Bytes")

	// Print network usage statistics
	for _, stat := range netStats {
		fmt.Fprintf(w, "%s\t%d\t%d\n", stat.Name, stat.BytesRecv, stat.BytesSent)
	}
}

//////////////////////////////////

func ExportToFile(result, filePath string) {

	if filePath != "" {
		// Replace newline characters with a space or other delimiter
		formattedData := strings.ReplaceAll(result, "\n", " ")

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logEntry := fmt.Sprintf("%s - %s\n", timestamp, formattedData)

		// Open the file in append mode, create it if it does not exist
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Write the log entry to the file
		if _, err := file.WriteString(logEntry); err != nil {
			fmt.Println("Error writing to file:", err)
		} else {
			fmt.Println("Output appended to", filePath)
		}
	}
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get CPU, Memory, Disk, Process, and Network usage",
	Run: func(cmd *cobra.Command, args []string) {
		cpuUsage := GetCPUUsage()
		memUsage := GetMemoryUsage()
		diskUsage := GetDiskUsage()

		// Format the output using tabwriter
		var formattedOutput bytes.Buffer
		w := tabwriter.NewWriter(&formattedOutput, 0, 0, 1, ' ', tabwriter.TabIndent)
		defer w.Flush()

		// Write the formatted system stats
		fmt.Fprintln(w, "System Stats:\t")
		fmt.Fprintln(w, "")

		// Write CPU usage
		fmt.Fprintf(w, "CPU Usage:\t%s\n", cpuUsage)
		fmt.Fprintln(w, "")

		// Write Memory usage
		fmt.Fprintf(w, "Memory Usage:\t%s\n", memUsage)
		fmt.Fprintln(w, "")

		// Write Disk usage
		fmt.Fprintf(w, "Disk Usage:\t%s\n", diskUsage)
		fmt.Fprintln(w, "")

		// Print the formatted output
		fmt.Println(formattedOutput.String())

		// Export to file if exportFilePath is provided
		if exportFilePath != "" {
			ExportToFile(formattedOutput.String(), exportFilePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringVarP(&exportFilePath, "export", "e", "", "Export to file (provide file path)")
}
