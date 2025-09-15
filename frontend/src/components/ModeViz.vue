<script lang="ts" setup>
import { reactive, ref, onMounted, computed, watch, nextTick } from 'vue'
import { main, viz } from '../../wailsjs/go/models';
import { useProjectStore } from '../project';
import * as THREE from 'three'
import chroma from 'chroma-js';

const project = useProjectStore()

const props = defineProps<{
    ModeData: viz.ModeData
    showNodePaths: boolean
}>()

watch(
    () => props.ModeData,
    (md) => { createFrames(md) },
    { deep: true }
)

watch(() => props.showNodePaths, (snp) => (nodePaths.visible = snp))

const showNodePaths = ref(true)

let frames = new Array<THREE.Group>;
let nodePaths = new THREE.Group;
let frameNum = 0;
let capturing = false;
let frameCenter = new THREE.Vector3;
let frameSize = new THREE.Vector3;
let clock = new THREE.Clock();
let delta = 0
const FOV = 10

// 2D Chart variables
const chartCanvas = ref<HTMLCanvasElement>()
let chartCtx: CanvasRenderingContext2D
let chartData = reactive({
    datasets: [
        { name: 'sin(x)', color: '#1f77b4', data: [] as {x: number, y: number}[] },
        { name: 'sin(2x)', color: '#ff7f0e', data: [] as {x: number, y: number}[] },
        { name: 'sin(3x)', color: '#2ca02c', data: [] as {x: number, y: number}[] },
        { name: 'sin(1.6x)', color: '#d62728', data: [] as {x: number, y: number}[] }
    ]
})

// Chart dimensions and margins
const chartMargins = { top: 40, right: 120, bottom: 60, left: 60 }

function generateChartData() {
    const configs = [
        { frequency: 0.05, amplitude: 70, phase: 0 },
        { frequency: 0.1, amplitude: 50, phase: 0 },
        { frequency: 0.15, amplitude: 40, phase: 0 },
        { frequency: 0.08, amplitude: 30, phase: Math.PI/4 }
    ]

    configs.forEach((config, index) => {
        const data: {x: number, y: number}[] = []
        for (let i = 0; i < 200; i++) {
            const x = (i - 100) * 0.1
            const y = config.amplitude * Math.sin(x * (config.frequency / 0.05) + config.phase)
            data.push({ x, y })
        }
        chartData.datasets[index].data = data
    })
}

function resizeChart() {
    if (!chartCanvas.value) return
    
    const container = chartCanvas.value.parentElement
    if (container) {
        chartCanvas.value.width = container.clientWidth
        chartCanvas.value.height = container.clientHeight
        drawChart()
    }
}

function drawChart() {
    if (!chartCtx || !chartCanvas.value) return

    const canvas = chartCanvas.value
    const ctx = chartCtx
    
    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    
    // Calculate chart area
    const chartWidth = canvas.width - chartMargins.left - chartMargins.right
    const chartHeight = canvas.height - chartMargins.top - chartMargins.bottom
    
    // Calculate data bounds
    let minX = Infinity, maxX = -Infinity, minY = Infinity, maxY = -Infinity
    chartData.datasets.forEach(dataset => {
        dataset.data.forEach(point => {
            minX = Math.min(minX, point.x)
            maxX = Math.max(maxX, point.x)
            minY = Math.min(minY, point.y)
            maxY = Math.max(maxY, point.y)
        })
    })

    // Helper functions
    const xScale = (x: number) => chartMargins.left + ((x - minX) / (maxX - minX)) * chartWidth
    const yScale = (y: number) => chartMargins.top + chartHeight - ((y - minY) / (maxY - minY)) * chartHeight

    // Draw background
    ctx.fillStyle = '#ffffff'
    ctx.fillRect(chartMargins.left, chartMargins.top, chartWidth, chartHeight)

    // Draw grid
    ctx.strokeStyle = '#e0e0e0'
    ctx.lineWidth = 1
    
    // Vertical grid lines
    for (let i = 0; i <= 10; i++) {
        const x = chartMargins.left + (i / 10) * chartWidth
        ctx.beginPath()
        ctx.moveTo(x, chartMargins.top)
        ctx.lineTo(x, chartMargins.top + chartHeight)
        ctx.stroke()
    }
    
    // Horizontal grid lines
    for (let i = 0; i <= 10; i++) {
        const y = chartMargins.top + (i / 10) * chartHeight
        ctx.beginPath()
        ctx.moveTo(chartMargins.left, y)
        ctx.lineTo(chartMargins.left + chartWidth, y)
        ctx.stroke()
    }

    // Draw axes
    ctx.strokeStyle = '#000000'
    ctx.lineWidth = 2
    
    // X-axis
    ctx.beginPath()
    ctx.moveTo(chartMargins.left, chartMargins.top + chartHeight)
    ctx.lineTo(chartMargins.left + chartWidth, chartMargins.top + chartHeight)
    ctx.stroke()
    
    // Y-axis
    ctx.beginPath()
    ctx.moveTo(chartMargins.left, chartMargins.top)
    ctx.lineTo(chartMargins.left, chartMargins.top + chartHeight)
    ctx.stroke()

    // Draw axis labels
    ctx.fillStyle = '#000000'
    ctx.font = '12px Arial'
    ctx.textAlign = 'center'
    
    // X-axis labels
    for (let i = 0; i <= 5; i++) {
        const x = chartMargins.left + (i / 5) * chartWidth
        const value = minX + (i / 5) * (maxX - minX)
        ctx.fillText(value.toFixed(1), x, chartMargins.top + chartHeight + 20)
    }
    
    // Y-axis labels
    ctx.textAlign = 'right'
    for (let i = 0; i <= 5; i++) {
        const y = chartMargins.top + chartHeight - (i / 5) * chartHeight
        const value = minY + (i / 5) * (maxY - minY)
        ctx.fillText(value.toFixed(0), chartMargins.left - 10, y + 4)
    }

    // Draw data lines
    chartData.datasets.forEach((dataset, datasetIndex) => {
        ctx.strokeStyle = dataset.color
        ctx.lineWidth = 2
        ctx.beginPath()
        
        let first = true
        dataset.data.forEach((point, pointIndex) => {
            const x = xScale(point.x)
            const y = yScale(point.y)
            
            if (x >= chartMargins.left && x <= chartMargins.left + chartWidth) {
                if (first) {
                    ctx.moveTo(x, y)
                    first = false
                } else {
                    ctx.lineTo(x, y)
                }
            }
        })
        ctx.stroke()
    })

    // Draw legend
    const legendX = chartMargins.left + chartWidth + 20
    chartData.datasets.forEach((dataset, index) => {
        const y = chartMargins.top + 20 + index * 25
        
        // Legend line
        ctx.strokeStyle = dataset.color
        ctx.lineWidth = 3
        ctx.beginPath()
        ctx.moveTo(legendX, y)
        ctx.lineTo(legendX + 20, y)
        ctx.stroke()
        
        // Legend text
        ctx.fillStyle = '#000000'
        ctx.font = '12px Arial'
        ctx.textAlign = 'left'
        ctx.fillText(dataset.name, legendX + 25, y + 4)
    })

    // Draw title
    ctx.fillStyle = '#000000'
    ctx.font = 'bold 16px Arial'
    ctx.textAlign = 'center'
    ctx.fillText('Sine Wave Comparison', canvas.width / 2, 25)

    // Draw axis titles
    ctx.font = '14px Arial'
    ctx.fillText('Time (s)', canvas.width / 2, canvas.height - 10)
    
    ctx.save()
    ctx.translate(15, canvas.height / 2)
    ctx.rotate(-Math.PI / 2)
    ctx.fillText('Amplitude', 0, 0)
    ctx.restore()
}

function create2DPlot() {
    generateChartData()
    nextTick(() => {
        if (chartCanvas.value) {
            resizeChart()
        }
    })
}

function createFrames(modeData: viz.ModeData) {
    if (modeData.Frames == null) return
    scene.clear()
    frames = [] as THREE.Group[];
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute([0, 0, 0], 3));
    const material = new THREE.PointsMaterial({ color: 0x888888 });
    const origin = new THREE.Points(geometry, material);
    origin.visible = false
    const allFramesGroup = new THREE.Group()
    allFramesGroup.add(origin)
    for (const f of modeData.Frames) {
        const frameGroup = new THREE.Group()
        for (const c of Object.values(f.Components)) {
            const curve = new THREE.CatmullRomCurve3(
                c.Line.map((p) => new THREE.Vector3(p.XYZ[0], p.XYZ[1], p.XYZ[2])))
            const points = curve.getPoints(50);
            const geometry = new THREE.BufferGeometry().setFromPoints(points);
            const material = new THREE.LineBasicMaterial({ color: 0xffffff, linewidth: 1 });
            const curveObject = new THREE.Line(geometry, material);
            frameGroup.add(curveObject)
            allFramesGroup.add(curveObject.clone()) // Add clone of object to be used for view sizing
        }
        frameGroup.visible = false // Initialize each group to not visible for animation
        frames.push(frameGroup)
        scene.add(frameGroup)
    }
    // Node paths
    const componentNames = Object.keys(modeData.Frames[0].Components)
    const curves = new Array<THREE.CatmullRomCurve3>
    const curveLengths = new Array<number>
    for (const compName of componentNames) {
        const numNodes = modeData.Frames[0].Components[compName].Line.length
        for (let i = 0; i < numNodes; i++) {
            let vectors = [] as THREE.Vector3[]
            for (const f of modeData.Frames) {
                const line = f.Components[compName].Line
                const p = line[i]
                vectors.push(new THREE.Vector3(p.XYZ[0], p.XYZ[1], p.XYZ[2]))
            }
            const curve = new THREE.CatmullRomCurve3(vectors)
            curve.closed = true
            curveLengths.push(curve.getLength())
            curves.push(curve)
        }
    }
    nodePaths.clear()
    const cs = chroma.scale(['008ae5', 'yellow']).domain([Math.min(...curveLengths), Math.max(...curveLengths)])
    for (let i = 0; i < curves.length; i++) {
        const points = curves[i].getPoints(50);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);
        const material = new THREE.LineBasicMaterial({ color: cs(curveLengths[i]).hex(), linewidth: 1, transparent: true })
        material.opacity = 0.8
        const line = new THREE.Line(geometry, material)
        line.computeLineDistances()
        nodePaths.add(line)
    }
    scene.add(nodePaths)
    const bb = new THREE.Box3().setFromObject(allFramesGroup);
    frameCenter = bb.getCenter(new THREE.Vector3())
    frameSize = bb.getSize(new THREE.Vector3())
    frameNum = 0
    const axesHelper = new THREE.AxesHelper(frameSize.x / 2);
    scene.add(axesHelper)

    // Create 2D plot
    create2DPlot();
}

let scene: THREE.Scene;
let renderer: THREE.WebGLRenderer;

const views = [
    {
        // Top View
        left: 0,
        bottom: 0.75,
        width: 0.4,
        height: 0.25,
        up: [1, 0, 0],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along Z axis to fit model in frame horizontally
            const fov = camera.fov * (Math.PI / 180);
            const fovh = 2 * Math.atan(Math.tan(fov / 2) * camera.aspect);
            let distance = 1.05 * (frameSize.y / 2 / Math.tan(fovh / 2) + frameSize.z)
            camera.position.fromArray([0, 0, distance]); // Looking along -Z (downward)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },
    {
        // Front View
        left: 0,
        bottom: 0,
        width: 0.4,
        height: 0.748,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along -X axis to fit model in frame vertically
            // See https://wejn.org/2020/12/cracking-the-threejs-object-fitting-nut/ for equation
            let distance = 1.15 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.x / 2)
            camera.position.fromArray([-distance, 0, frameCenter.z]); // Looking along X (downwind)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },
    {
        // Side View
        left: 0.402,
        bottom: 0,
        width: 0.25,
        height: 0.748,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along -Y axis to fit model in frame vertically
            let distance = 1.15 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.y / 2)
            camera.position.fromArray([0, -distance, frameCenter.z]); // Looking along -Y (side)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },

    {
        // Isometric View
        left: 0.654,
        bottom: 0,
        width: 0.35,
        height: 1.0,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along Z axis to fit model in frame horizontally
            let distanceFront = 1.0 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.x / 2)
            let distanceSide = 1.0 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.y / 2)
            camera.position.fromArray([-distanceFront, -distanceSide, frameCenter.z + 3 * frameSize.z]); // Looking along -Z (downward)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    }
];

function animate() {
    requestAnimationFrame(animate);
    delta += clock.getDelta()
    if (delta > 1.5 / frames.length) {
        delta = 0
        if (frames.length > 0) {
            frames[frameNum].visible = false;
            frameNum++
            if (frameNum >= frames.length) frameNum = 0
            frames[frameNum].visible = true;
        }
        render();
    }
}

function render() {
    const canvas = renderer.domElement;
    const canvasWidth = canvas.clientWidth;
    const canvasHeight = canvas.clientHeight;
    if (canvas.width !== canvasWidth || canvas.height !== canvasHeight) {
        renderer.setSize(canvasWidth, canvasHeight, false);
    }

    for (let ii = 0; ii < views.length; ++ii) {
        const view = views[ii];
        const camera = view.camera;

        view.updateCamera(camera);

        const left = Math.floor(canvasWidth * view.left);
        const bottom = Math.floor(canvasHeight * view.bottom);
        const width = Math.floor(canvasWidth * view.width);
        const height = Math.floor(canvasHeight * view.height);

        renderer.setViewport(left, bottom, width, height);
        renderer.setScissor(left, bottom, width, height);
        renderer.setScissorTest(true);

        camera.aspect = width / height;
        camera.updateProjectionMatrix();

        renderer.render(scene, camera);
    }
}

onMounted(() => {
    const canvas = <HTMLCanvasElement>document.getElementById('modeVizCanvas')!;

    for (let ii = 0; ii < views.length; ++ii) {
        const view = views[ii];
        const camera = new THREE.PerspectiveCamera(FOV, 2, 1, 10000);
        camera.up.fromArray(view.up);
        view.updateCamera(camera)
        view.camera = camera;
    }

    scene = new THREE.Scene();
    createFrames(props.ModeData)

    renderer = new THREE.WebGLRenderer({ antialias: true, canvas });
    renderer.setClearColor(0x3a3b3c);

    // Initialize chart
    if (chartCanvas.value) {
        chartCtx = chartCanvas.value.getContext('2d')!
        resizeChart()
    }

    // Handle window resize
    window.addEventListener('resize', resizeChart)

    animate();
})

</script>

<template>
    <div class="h-100 w-100 d-flex flex-column">
        <!-- 3D Mode Visualization - Takes 80% height -->
        <div class="flex-grow-1" style="height: 80%;">
            <canvas id="modeVizCanvas" class="h-100 w-100"></canvas>
        </div>
        
        <!-- 2D Chart - Takes 20% height -->
        <div class="border-top bg-white" style="height: 20%;">
            <canvas ref="chartCanvas" class="w-100 h-100"></canvas>
        </div>
    </div>
</template>