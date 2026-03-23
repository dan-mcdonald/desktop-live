<script lang="ts">
  import { PollMessage, PostMessage } from "../wailsjs/go/main/App.js";
  import { onMount } from "svelte";

  let resultText: string = "Please enter your name below 👇";
  let name: string;
  let recorder: MediaRecorder;
  let chunks: BlobPart[] = [];

  function onDataAvailable(e: BlobEvent) {
    chunks.push(e.data);
  }

  async function onReady(): Promise<void> {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    recorder = new MediaRecorder(stream);
    recorder.ondataavailable = onDataAvailable;
    recorder.onstop = onStop;
    // stream.getTracks().forEach((track) => track.stop());
    // const devices = await navigator.mediaDevices.enumerateDevices();
    // const inputStream = devices.find((device) => device.kind === "audioinput");
    // resultText = JSON.stringify(inputStream);
    const audioCtx = new AudioContext();
    const source = audioCtx.createMediaStreamSource(stream);
    const bufferLength = 2048;
    const analyser = audioCtx.createAnalyser();
    analyser.fftSize = bufferLength;
    const dataArray = new Uint8Array(bufferLength);
    source.connect(analyser);
    const canvas = document.getElementById("visualizer")! as HTMLCanvasElement;
    const canvasCtx = canvas.getContext("2d")!;
    draw();

    function draw() {
      const WIDTH = canvas.width;
      const HEIGHT = canvas.height;

      requestAnimationFrame(draw);

      analyser.getByteTimeDomainData(dataArray);

      canvasCtx.fillStyle = "rgb(200, 200, 200)";
      canvasCtx.fillRect(0, 0, WIDTH, HEIGHT);

      canvasCtx.lineWidth = 2;
      canvasCtx.strokeStyle = "rgb(0, 0, 0)";

      canvasCtx.beginPath();

      let sliceWidth = (WIDTH * 1.0) / bufferLength;
      let x = 0;

      for (let i = 0; i < bufferLength; i++) {
        let v = dataArray[i] / 128.0;
        let y = (v * HEIGHT) / 2;

        if (i === 0) {
          canvasCtx.moveTo(x, y);
        } else {
          canvasCtx.lineTo(x, y);
        }

        x += sliceWidth;
      }

      canvasCtx.lineTo(canvas.width, canvas.height / 2);
      canvasCtx.stroke();
    }
  }

  onMount(onReady);

  function startTalk() {
    chunks = [];
    recorder.start();
  }

  function stopTalk() {
    recorder.stop();
  }

  function onStop() {
    console.log(`Recording stopped, chunks: ${chunks.length}`);
    console.log("first chunk:", chunks[0]);
  }

  function greet(): void {
    PostMessage({ Text: name }).then((result) => {
      resultText = "";
    });
  }

  function pollMessage(): void {
    PollMessage().then((result) => {
      console.log(JSON.stringify(result));
      pollMessage();
    });
  }

  pollMessage();
</script>

<main>
  <canvas id="visualizer" width="500" height="200"></canvas>
  <button id="talk" on:pointerdown={startTalk} on:pointerup={stopTalk}
    >Talk</button
  >
  <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <input
      autocomplete="off"
      bind:value={name}
      class="input"
      id="name"
      type="text"
    />
    <button class="btn" on:click={greet}>Greet</button>
  </div>
</main>

<style>
  #visualizer {
    display: block;
    margin: auto;
  }
  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }
</style>
