// SPDX-License-Identifier: MIT
/// @author Perry
pragma solidity ^0.8.20;

import "./lib/Auth.sol";
import "./lib/Channel.sol";

contract ChannelManger {
    using Channel for Channel.MetaData;
    using Auth for Auth.Role;

    struct ChannelInfo {
        Channel.MetaData data;
        Auth.Role role;
    }

    mapping(uint => ChannelInfo) private channels;

    event ChannelCreated(uint indexed channelId, address indexed creator);
    event ChannelDeleted(uint indexed channelId);
    event JoinedChannel(uint indexed channelId, address indexed user);
    event LeftChannel(uint indexed channelId, address indexed user);
    event RoleAdded(
        uint indexed channelID,
        address indexed user,
        Auth.RoleType roleType
    );
    event RoleRemoved(
        uint indexed channelID,
        address indexed user,
        Auth.RoleType roleType
    );

    /// @notice create channel
    function createChannel(
        uint channelId,
        Channel.ChannelType channelType
    ) external {
        channels[channelId].data.create(msg.sender, channelType);
        channels[channelId].role.addRole(msg.sender, Auth.RoleType.Admin);

        emit ChannelCreated(channelId, msg.sender);
        emit RoleAdded(channelId, msg.sender, Auth.RoleType.Admin);
    }

    /// @notice delete channel
    function deleteChannel(uint channelId) external {
        require(
            channels[channelId].role.hasRole(msg.sender, Auth.RoleType.Admin)
        );
        channels[channelId].data.remove();
        emit ChannelDeleted(channelId);
    }

    /// @notice join channel
    function joinChannel(uint channelId) external {
        channels[channelId].data.join(msg.sender);
        channels[channelId].role.addRole(msg.sender, Auth.RoleType.Member);
        emit JoinedChannel(channelId, msg.sender);
        emit RoleAdded(channelId, msg.sender, Auth.RoleType.Member);
    }

    /// @notice leave channel
    function leaveChannel(uint channelId) external {
        channels[channelId].data.leave(msg.sender);
        channels[channelId].role.removeRole(msg.sender, Auth.RoleType.Member);
        emit LeftChannel(channelId, msg.sender);
        emit RoleRemoved(channelId, msg.sender, Auth.RoleType.Member);
    }
}
