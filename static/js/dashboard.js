// Initialize global state
const dashboardState = {
    currentTask: null,
    nextTask: null
};

// Cache DOM elements
const themeToggle = document.getElementById('theme-toggle');
const body = document.body;
const dashboardTitle = document.getElementById('dashboard-title');
const circleTimer = document.getElementById('circle-timer');
const digitalClock = document.getElementById('digital-clock');
const timerDisplay = document.getElementById('timer-display');
const progressCircle = document.getElementById('progress-circle');
const timeDisplay = document.getElementById('time-display');
const dateDisplay = document.getElementById('date-display');
const nextTaskName = document.getElementById('next-task-name');
const nextTaskTime = document.getElementById('next-task-time');
const nextTaskCountdown = document.getElementById('next-task-countdown');

// Theme toggle functionality
if (themeToggle) {
    themeToggle.addEventListener('click', () => {
        body.classList.toggle('dark-mode');
        localStorage.setItem('darkMode', body.classList.contains('dark-mode'));
    });
}

if (localStorage.getItem('darkMode') === 'true') {
    body.classList.add('dark-mode');
}

async function fetchDashboardData() {
    try {
        const response = await fetch('/api/dashboard-data');
        if (!response.ok) {
            throw new Error(`API error: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error("Error fetching dashboard data:", error);
        return { currentTask: null, nextTask: null };
    }
}

function formatTime(dateString) {
    return new Date(dateString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString([], { weekday: 'long', month: 'long', day: 'numeric' });
}

function getTimeComponents(timeLeft) {
    return {
        days: Math.floor(timeLeft / (1000 * 60 * 60 * 24)),
        hours: Math.floor((timeLeft % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
        minutes: Math.floor((timeLeft % (1000 * 60 * 60)) / (1000 * 60)),
        seconds: Math.floor((timeLeft % (1000 * 60)) / 1000)
    };
}

function formatCountdown(timeLeft) {
    if (timeLeft <= 0) return "00:00:00";
    const { hours, minutes, seconds } = getTimeComponents(timeLeft);
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

function formatNextTaskCountdown(timeLeft) {
    if (timeLeft <= 0) return "";
    const { days, hours, minutes } = getTimeComponents(timeLeft);
    if (days > 0) return `in ${days}d ${hours}h ${minutes}m`;
    if (hours > 0) return `in ${hours}h ${minutes}m`;
    return `in ${minutes}m`;
}

function updateCurrentTaskCircle(task) {
    const now = Date.now();
    const start = new Date(task.startDate).getTime();
    const end = new Date(task.endDate).getTime();
    const progress = Math.min(100, Math.max(0, ((now - start) / (end - start)) * 100));
    const timeLeft = end - now;

    circleTimer.style.display = 'block';
    digitalClock.style.display = 'none';

    dashboardTitle.textContent = task.name || 'Current Task';
    const circumference = 2 * Math.PI * 115;
    const offset = circumference * (1 - progress / 100);

    progressCircle.style.strokeDasharray = circumference;
    progressCircle.style.strokeDashoffset = offset;
    timerDisplay.textContent = formatCountdown(timeLeft);
}

function updateDigitalClock() {
    const now = new Date();

    circleTimer.style.display = 'none';
    digitalClock.style.display = 'block';

    dashboardTitle.textContent = 'Task Dashboard';
    timeDisplay.textContent = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    dateDisplay.textContent = now.toLocaleDateString([], { weekday: 'long', month: 'long', day: 'numeric' });
}

function updateNextTaskInfo(task) {
    if (!task) {
        nextTaskName.textContent = 'No upcoming tasks';
        nextTaskTime.textContent = '';
        nextTaskCountdown.textContent = '';
        return;
    }

    nextTaskName.textContent = task.name || 'Unnamed Task';
    nextTaskTime.textContent = formatTime(task.startDate);

    const now = Date.now();
    const start = new Date(task.startDate).getTime();
    nextTaskCountdown.textContent = formatNextTaskCountdown(start - now);
}

function updateActiveTimers() {
    if (!dashboardState.currentTask) {
        updateDigitalClock();
    } else {
        updateCurrentTaskCircle(dashboardState.currentTask);
    }
    updateNextTaskInfo(dashboardState.nextTask);
}

function startIntervals() {
    setInterval(() => {
        const now = Date.now();

        if (digitalClock.style.display === 'block') {
            updateDigitalClock();
        }

        if (circleTimer.style.display === 'block' && dashboardState.currentTask) {
            const startTime = new Date(dashboardState.currentTask.startDate).getTime();
            const endTime = new Date(dashboardState.currentTask.endDate).getTime();
            const totalDuration = endTime - startTime;
            const elapsed = now - startTime;
            const timeLeft = endTime - now;

            if (timeLeft > 0) {
                const progress = Math.min(100, Math.max(0, (elapsed / totalDuration) * 100));
                const circumference = 2 * Math.PI * 115;
                const offset = circumference * (1 - progress / 100);

                progressCircle.style.strokeDasharray = circumference;
                progressCircle.style.strokeDashoffset = offset;
                timerDisplay.textContent = formatCountdown(timeLeft);
            } else {
                refresh();
            }
        }

        if (dashboardState.nextTask) {
            const startTime = new Date(dashboardState.nextTask.startDate).getTime();
            const timeUntilStart = startTime - now;

            if (timeUntilStart > 0) {
                nextTaskCountdown.textContent = formatNextTaskCountdown(timeUntilStart);
            } else {
                refresh();
            }
        }
    }, 1000);
}

async function refresh() {
    try {
        const data = await fetchDashboardData();
        dashboardState.currentTask = data.currentTask;
        dashboardState.nextTask = data.nextTask;
        updateActiveTimers();
    } catch (error) {
        console.error("Error refreshing dashboard:", error);
    }
}

function init() {
    refresh();
    startIntervals();
    setInterval(refresh, 60000);
}

document.addEventListener('DOMContentLoaded', init);