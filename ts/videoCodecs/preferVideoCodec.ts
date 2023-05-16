import { hasAnyTexts } from "ts/hasAnyTexts";

/** Prefer Specified Video Codec */
export function preferVideoCodec(codecs: RTCRtpCodecCapability[], mimeType: string): RTCRtpCodecCapability[] {
    // skip sorting video codecs if the mimeType is empty.
    if(!hasAnyTexts(mimeType)) {
        return codecs;
    }
    const otherCodecs: RTCRtpCodecCapability[] = [];
    const sortedCodecs: RTCRtpCodecCapability[] = [];
    codecs.forEach((codec) => {
      if (codec.mimeType === mimeType) {
        sortedCodecs.push(codec);
      } else {
        otherCodecs.push(codec);
      }
    });  
    return sortedCodecs.concat(otherCodecs);
  }