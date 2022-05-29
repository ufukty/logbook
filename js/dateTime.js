export function dateToEpoch(dateTimeString) {
    var date = new Date(dateTimeString.substring(0, 10));
    return date.valueOf().toString();
}

export function classifyTasksByDays(tasks) {
    var classifiedTasks = {}; // type: { "date1": [task1, task2], "date2": [...], ...}
    for (const taskId in tasks) {
        if (Object.hasOwnProperty.call(tasks, taskId)) {
            var completionDayEpoch = dateToEpoch(tasks[taskId].createdAt);
            if (!classifiedTasks.hasOwnProperty(completionDayEpoch)) {
                classifiedTasks[completionDayEpoch] = [];
            }
            classifiedTasks[completionDayEpoch].push(taskId);
        }
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
