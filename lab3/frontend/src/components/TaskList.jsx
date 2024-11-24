// src/components/TaskList.js
import React, { useEffect, useState } from 'react';
import { useTaskContext } from '../context/TaskContext';

const TaskList = () => {
    const { token } = useTaskContext();
    const [tasks, setTasks] = useState([]);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const response = await fetch('http://localhost:8080/tasks', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`, 
                    },
                });

                if (!response.ok) {
                    throw new Error('Unauthorized');
                }

                const result = await response.json();
                setTasks(result); 
            } catch (err) {
                setError(err.message);
            }
        };

        if (token) {
            fetchTasks();
        }
    }, [token]);

    return (
        <div>
            {error && <p>{error}</p>}
            <ul>
                {tasks.map(task => (
                    <li key={task.id}>{task.title} - {task.description}</li>
                ))}
            </ul>
        </div>
    );
};

export default TaskList;
