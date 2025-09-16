// SPDX-License-Identifier: MIT
/// @author Perry
pragma solidity ^0.8.20;
library Channel {
    enum ChannelType {Public, Private};

    struct MetaData {
        bool exists;
        ChannelType channelType;
        address creator;
        mapping(address => bool) members;
        address[] memberList;
    }

    error ChannelAlreadyExists();
    error ChannelNotExist();
    error AlreadyChannelMember();
    error NotChannelMember();

    /// @dev
    function create(MetaData storage ch, address creator, ChannelType channelType) internal {
        if (ch.exists) revert ChannelAlreadyExists();
        ch.exists = true;
        ch.creator = creator;
        ch.channelType = channelType;
        ch.members[creator] = true;
        ch.memberList.push(creator);
    }

    /// @dev
    function remove(MetaData storage ch) internal {
        if (!ch.exists) revert ChannelNotExist();
        ch.exists = false;
    }

    /// @dev
    function join(MetaData storage ch, address user) internal {
        if (!ch.exists) revert ChannelNotExist();
        if (ch.members[user]) revert AlreadyChannelMember();
        ch.members[user] = true;
        ch.memberList.push(user);
    }

    /// @dev
    function leave(MetaData storage ch, address user) internal {
        if (!ch.exists) revert ChannelNotExist();
        if (!ch.members[user]) revert NotChannelMember();
        ch.members[user] = false;
        uint len = memberList.length;
        for (uint i = 0; i < len; i++) {
            if (memberList[i] == user) {
                memberList[i] == memberList[len - 1]
                memberList.pop();
                break;
            }
        }
    }

    /// @dev
    function isChannelMember(MetaData storage ch, address user) internal view returns (bool) {
        return ch.members[user];
    }

    /// @dev
    function getChannelMembers(MetaData storage ch) internal view returns (address[] memory) {
        return ch.memberList;
    }
}