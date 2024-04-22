import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: "0.8.24",
  networks: {
    hardhat: {
    },
    auradev: {
      url: "https://jsonrpc.dev.aura.network",
    }
  }
};

export default config;
