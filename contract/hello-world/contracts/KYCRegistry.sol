pragma solidity >=0.7.3;

contract KYCRegistry {
    mapping (address => bool) public kycedAddresses;

    function isKYCed(address _address) public view returns (bool) {
        return kycedAddresses[_address];
    }

    function setKYCStatus(address _address, bool _status) public {
        kycedAddresses[_address] = _status;
    }
}
