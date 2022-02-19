export function dateToEpoch(dateTimeString) {
    var date = new Date(dateTimeString.substring(0, 10));
    return date.valueOf().toString();
}

export function classifyTasksByDays(linear_tasks) {
    var classifiedTasks = {}; // type: { "date1": [task1, task2], "date2": [...], ...}
    for (const task of linear_tasks) {
        var backendTimestamp = task.completed_at;
        var completionDayEpoch = dateToEpoch(backendTimestamp);
        if (!classifiedTasks.hasOwnProperty(completionDayEpoch)) {
            classifiedTasks[completionDayEpoch] = [];
        }
        classifiedTasks[completionDayEpoch].push(task);
    }
    return classifiedTasks;
}

var intl = new Intl.DateTimeFormat("default", {
    // weekday: "long",
    year: "2-digit",
    month: "long",
    day: "numeric",
});

export function timestampToLocalizedText(date) {
    // var date = new Date(backendTimestamp);
    var completionDay = intl.format(date);
    return completionDay;
}
