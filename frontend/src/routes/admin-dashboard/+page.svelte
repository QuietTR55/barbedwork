<script lang="ts">
    import { onMount } from 'svelte';
	import {
		Chart,
		LineController,
		LineElement,
		PointElement,
		LinearScale,
		CategoryScale,
		Tooltip,
		Legend,
		DoughnutController,
		ArcElement
	} from 'chart.js';
	import Icon from '@iconify/svelte';
    import type { PageData } from './$types';

    let { data }: { data: PageData } = $props();

    let chartCanvas = $state<HTMLCanvasElement>();

	Chart.register(
		LineController,
		LineElement,
		PointElement,
		LinearScale,
		CategoryScale,
		Tooltip,
		Legend,
		DoughnutController,
		ArcElement
	);

    onMount(() => {
        const ctx = chartCanvas?.getContext('2d');
        if (ctx) {
            new Chart(ctx, {
                type: 'doughnut',
                data: {
                    labels: ['Pending Tasks', 'Completed Tasks', 'In Progress Tasks'],
                    datasets: [{
                        label: 'Task Status',
                        data: [10, 25, 5], // Example data
                        backgroundColor: [
                            '#8b5cf6', // Purple for Pending
                            '#10b981', // Green for Completed
                            '#3b82f6'  // Blue for In Progress
                        ],
                        hoverOffset: 1
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                            position: 'bottom',
                            align: 'center',
                            labels: {
                                usePointStyle: true,
                                pointStyle: 'circle',
                                padding: 10,
                                font: {
                                    size: 12,
                                    weight: 'bold'
                                }
                            }
                        }
                    }
                }
            });
        }
    });
</script>
<div class="flex flex-col lg:flex-wrap md:flex-wrap md:flex-row gap-4 p-4 overflow-y-auto h-full w-full">
    <div class="stats-card stats-card-admin-dashboard-sizing">
        <p class="text-sm font-medium">Tasks</p>
        <div class="w-[80%] h-[80%] p-0 m-auto">
            <canvas bind:this={chartCanvas} id="myChart"></canvas>
        </div>
    </div>
    <div class="stats-card stats-card-admin-dashboard-sizing">
        <p class="text-lg font-medium">Member Count</p>
        <p class="text-4xl font-bold text-accent">100</p>
    </div>
</div>