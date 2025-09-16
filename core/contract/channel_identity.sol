// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ChannelAuthManger {
    struct Channel {
        mapping(address => bool) public admins;
        mapping(address => bool) public members;
        bool exists;
    }

    mapping(uint => Channel) private channels;

    event ChannelCreated(uint indexed channelId, address indexed creator);
    event AdminAdded(uint indexed channelId, address indexed admin);
    event MemberAdded(uint indexed channelId, address indexed member);

    function createChannel(uint channelId) external {
        require(!channels[channelId].exists, "Channel already exists");
        Channel storage ch = channels[channelId];
        ch.exists = true;
        ch.admins[msg.sender] = true;
        emit ChannelCreated(channelId, msg.sender);
    }

    modifier onlyAdmin(uint channelId) {
        require(channels[channelId].admins[msg.sender], "Not admin");
        _;
    }

    function addAdmin(uint channelId, address _admin) external onlyAdmin(channelId) {
        require(!channels[channelId].admins[_admin], "Already an admin.");
        channels[channelId].admins[_admin] = true;
        emit AdminAdded(channelId, _admin);
    }

    function addMember(uint channelId, address _member) external onlyAdmin(channelId) {
        require(!channels[channelId].members[_member], "Already a member.");
        channels[channelId].members[_member] = true;
        emit MemberAdded(channelId,_member);
    }

    function isAdmin(uint channelId, address _addr) external view returns (bool) {
        return channels[channelId].admins[_addr];
    }

    function isMember(uint channelId, address _addr) external view returns (bool) {
        return channels[channelId].members[_addr];
    }
}
