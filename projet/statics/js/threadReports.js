function ResolveReport(threadName, reportId) {
    setReportToResolved(threadName, reportId)
        .then(r => {
            if (r.ok) {
                // If the report was resolved successfully, remove the report from the list
                document.getElementById(`report-${reportId}`).remove();
                // Show a success message
                alert('Report ' + reportId + ' has been resolved.');
            }
        }).catch(error => {
            alert('Error resolving report: ' + error);
            console.error("Error:", error);
        });
}