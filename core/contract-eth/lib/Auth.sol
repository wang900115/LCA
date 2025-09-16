// SPDX-License-Identifier: MIT
/// @author Perry
pragma solidity ^0.8.20;
library Auth {
    enum  RoleType {
        Admin,
        Member,
        // Candidate,
        // Validator,
    }
    struct Role {
        mapping(address => bool) admins;
        mapping(address => bool) members;
        // mapping(address => bool) candidates;
        // mapping(address => bool) validators;
    }

    error NotAdmin();
    error NotMember();
    // error NotValidator();
    // error NotCandidate();
    error AlreadyAdmin();
    error AlreadyMember();
    // error AlreadyValidator();
    // error AlreadyCandidate();

    /// @dev add
    function addRole(Role storage role, address user, RoleType roleType) internal {
        if (roleType == RoleType.Admin) {
            if (role.admins[user]) revert AlreadyAdmin();
            role.admins[user] = true;
        } else if (roleType == RoleType.Member) {
            if (role.members[user]) revert AlreadyMember();
            role.members[user] = true;
        }
        // } else if (roleType == RoleType.Candidate) {
        //     if (role.candidates[user]) revert AlreadyCandidate();
        //     role.candidates[user] = true;
        // } else if (roleType == RoleType.Validator) {
        //     if (role.validators[user]) revert AlreadyValidator();
        //     role.validators[user] = true;
        // }
    }

    /// @dev remove
    function removeRole(Role storage role, address user, RoleType roleType) internal {
        if (roleType == RoleType.Admin) {
            if (!role.admins[user]) revert NotAdmin();
            role.admins[user] = false;
        } else if (roleType == RoleType.Member) {
            if (!role.members[user]) revert NotMember();
            role.members[user] = false;
        }
        // } else if (roleType == RoleType.Candidate) {
        //     if (!role.candidates[user]) revert NotCandidate();
        //     role.candidates[user] = false;
        // } else if (roleType == RoleType.Validator) {
        //     if (!role.validators[user]) revert NotValidator();
        //     role.validators[user] = false;
        // }
    }

    /// @dev exist
    function hasRole(Role storage role, address user, RoleType roleType) internal view returns (bool) {
        if (roleType == RoleType.Admin) return role.admins[user];
        if (roleType == RoleType.Member) return role.members[user];
        // if (roleType == RoleType.Candidate) return role.candidates[user];
        // if (roleType == RoleType.Validator) return role.validators[user];
        return false;
    }
}
