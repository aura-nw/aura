import { ETH } from '@evmos/address-converter';
import { fromBech32, toBech32 } from '@cosmjs/encoding';

function makeBech32Encoder(prefix: string) {
  return (data: Uint8Array) => toBech32(prefix, data);
}

function makeBech32Decoder(currentPrefix: string) {
  return (input: string) => {
    const { prefix, data } = fromBech32(input);
    if (prefix !== currentPrefix) {
      throw Error('Unrecognised address format');
    }
    return Buffer.from(data);
  };
}

export function convertBech32AddressToEthAddress(
  prefix: string,
  bech32Address: string
) {
  const data = makeBech32Decoder(prefix)(bech32Address);
  return ETH.encoder(data);
}

export function convertEthAddressToBech32Address(
  prefix: string,
  ethAddress: string
) {
  const data = ETH.decoder(ethAddress);
  return makeBech32Encoder(prefix)(data);
}