// This code is based on GoogleChromeLabs/web-audio-samples(Copyright (c) 2022 The Chromium Authors) for reference
// https://github.com/GoogleChromeLabs/web-audio-samples

/* global currentTime */

const FRAME_INTERVAL = 1 / 60;

/**
 *  Measure microphone volume.
 *
 * @class VolumeMeter
 * @extends AudioWorkletProcessor
 */
class VolumeMeasurer extends AudioWorkletProcessor {

  constructor() {
    super();
    this._lastUpdate = currentTime;
  }

  calculateRMS(inputChannelData) {
    // Calculate the squared-sum.
    let sum = 0;
    // the value of "inputChannelData.length" is 128 by default.
    for (let i = 0; i < inputChannelData.length; i++) {
      sum += inputChannelData[i] * inputChannelData[i];
    }
    // Calculate the RMS(Root Mean Square) level.
    return Math.sqrt(sum / inputChannelData.length);
  }
  // "output" and "parameters" can be omitted
  process(inputs) {
    // This example only handles mono channel.
    const inputChannelData = inputs[0][0];
    // Calculate and post the RMS level every 16ms.
    if (currentTime - this._lastUpdate > FRAME_INTERVAL) {
      const volume = this.calculateRMS(inputChannelData);
      this.port.postMessage(volume);
      this._lastUpdate = currentTime;
    }
    return true;
  }
}

registerProcessor("volume-measurer", VolumeMeasurer);