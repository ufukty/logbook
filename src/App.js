import TaskDay from "./ui-components/task-day/TaskDay";

import "./App.css";

function App() {
    return (
        <div className="document-sheet">
            <a
                id="home-button"
                class="floating-corner left top"
                href="index.html"
            >
                Logbook
            </a>

            <div id="sheet-settings" class="floating-corner right top dark">
                <div>Share</div>

                <div>Sync</div>
            </div>

            <div id="active-task-details" class="floating-corner left bottom">
                History for active task
                <div class="task">PAM for SSH</div>
                <div class="task">ACL - Redis</div>
                <div class="task">TOTP for SSH</div>
            </div>

            <div id="date-anchors" class="floating-corner right bottom">
                <a href="#august-13-2021">13th August</a>
                <a href="#august-14-2021">14th August</a>
            </div>
            <TaskDay />
        </div>
    );
}

export default App;
