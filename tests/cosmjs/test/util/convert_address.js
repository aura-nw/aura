import { ETH } from '@evmos/address-converter';
import { fromBech32, toBech32 } from '@cosmjs/encoding';

function makeBech32Encoder(prefix) {
  return (data) => toBech32(prefix, data);
}

function makeBech32Decoder(currentPrefix) {
  return (input) => {
    const { prefix, data } = fromBech32(input);
    if (prefix !== currentPrefix) {
      throw Error('Unrecognised address format');
    }
    return Buffer.from(data);
  };
}

export function convertBech32AddressToEthAddress(
  prefix,
  bech32Address
) {
  const data = makeBech32Decoder(prefix)(bech32Address);
  return ETH.encoder(data);
}

export function convertEthAddressToBech32Address(
  prefix,
  ethAddress
) {
  const data = ETH.decoder(ethAddress);
  return makeBech32Encoder(prefix)(data);
}