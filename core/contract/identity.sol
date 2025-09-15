// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ChatAuth {
    mapping(address => bool) public admins;
    mapping(address => bool) public members;

    event AdminAdded(address indexed admin);
    event MemberAdded(address indexed member);

    constructor() {
        admins[msg.sender] = true;
        emit AdminAdded(msg.sender);
    }

    modifier onlyAdmin() {
        require(admins[msg.sender], "Not admin");
        _;
    }

    function addAdmin(address _admin) external onlyAdmin {
        admins[_admin] = true;
        emit AdminAdded(_admin)
    }

    function addMember(address _member) external onlyAdmin {
        members[_member] = true;
        emit MemberAdded(_member)
    }

    function isAdmin(address _addr) external view returns (bool) {
        return admins[_addr];
    }

    function isMember(address _addr) external view returns (bool) {
        return members[_addr];
    }
}
