async function main() {
   const KYCRegistry = await ethers.getContractFactory("KYCRegistry");

   // Start deployment, returning a promise that resolves to a contract object
   const kyc_registry = await KYCRegistry.deploy();
   console.log("Contract deployed to address:", kyc_registry.address);
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error);
    process.exit(1);
  });


