import React, { useState, useEffect } from 'react';
import { Plus, X } from 'lucide-react';

interface Task {
  id: string;
  title: string;
  description: string;
  status: 'todo' | 'in-progress' | 'completed';
  assignee: string;
  dueDate: string;
}

const API_BASE_URL = 'http://localhost:8080/api';

const TaskManager: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [showNewTask, setShowNewTask] = useState(false);
  const [newTask, setNewTask] = useState<Partial<Task>>({
    title: '',
    description: '',
    status: 'todo',
    assignee: '',
    dueDate: '',
  });

  // ✅ Fetch tasks from backend (Uses cookies instead of headers)
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/gettasks`, {
          credentials: 'include', // ✅ Ensure cookies are sent with the request
        });

        if (!response.ok) throw new Error('Failed to fetch tasks');

        const data: Task[] = await response.json();
        setTasks(data);
      } catch (error) {
        console.error(error);
      }
    };

    fetchTasks();
  }, []);

  // ✅ Create a new task
  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await fetch(`${API_BASE_URL}/tasks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include', // ✅ Uses cookies
        body: JSON.stringify(newTask),
      });

      if (!response.ok) throw new Error('Failed to create task');

      const createdTask = await response.json();
      setTasks([...tasks, createdTask]);
      setShowNewTask(false);
      setNewTask({ title: '', description: '', status: 'todo', assignee: '', dueDate: '' });
    } catch (error) {
      console.error(error);
    }
  };

  // ✅ Update task status
  const handleStatusChange = async (taskId: string, newStatus: Task['status']) => {
    try {
      await fetch(`${API_BASE_URL}/tasks/${taskId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include', // ✅ Uses cookies
        body: JSON.stringify({ status: newStatus }),
      });

      setTasks(tasks.map(task => (task.id === taskId ? { ...task, status: newStatus } : task)));
    } catch (error) {
      console.error(error);
    }
  };

  // ✅ Delete task
  const handleDeleteTask = async (taskId: string) => {
    try {
      await fetch(`${API_BASE_URL}/tasks/${taskId}`, {
        method: 'DELETE',
        credentials: 'include', // ✅ Uses cookies
      });

      setTasks(tasks.filter(task => task.id !== taskId));
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Task Manager</h1>
        <button onClick={() => setShowNewTask(true)} className="bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded-lg font-medium flex items-center gap-2">
          <Plus size={20} /> New Task
        </button>
      </div>

      {/* Task Board */}
      <div className="grid grid-cols-3 gap-6">
        {['todo', 'in-progress', 'completed'].map(status => (
          <TaskColumn key={status} title={status} tasks={tasks.filter(task => task.status === status)} onStatusChange={handleStatusChange} onDelete={handleDeleteTask} />
        ))}
      </div>
    </div>
  );
};

interface TaskColumnProps {
  title: string;
  tasks: Task[];
  onStatusChange: (taskId: string, status: Task['status']) => void;
  onDelete: (taskId: string) => void;
}

const TaskColumn: React.FC<TaskColumnProps> = ({ title, tasks, onStatusChange, onDelete }) => (
  <div className="bg-gray-800 rounded-lg p-4">
    <h2 className="text-xl font-semibold mb-4">{title}</h2>
    {tasks.map(task => (
      <TaskCard key={task.id} task={task} onStatusChange={onStatusChange} onDelete={onDelete} />
    ))}
  </div>
);

interface TaskCardProps {
  task: Task;
  onStatusChange: (taskId: string, status: Task['status']) => void;
  onDelete: (taskId: string) => void;
}

const TaskCard: React.FC<TaskCardProps> = ({ task, onStatusChange, onDelete }) => (
  <div className="bg-gray-700 rounded-lg p-4">
    <h3 className="font-medium">{task.title}</h3>
    <p className="text-sm text-gray-400 mb-3">{task.description}</p>
    <div className="flex justify-between items-center text-sm">
      <span className="text-gray-400">Due: {task.dueDate}</span>
      <button onClick={() => onDelete(task.id)} className="text-red-500 hover:text-red-700">
        <X size={16} />
      </button>
    </div>
  </div>
);

export default TaskManager;
