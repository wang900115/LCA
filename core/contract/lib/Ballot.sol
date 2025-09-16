// SPDX-License-Identifier: MIT
/// @author Perry
pragma solidity ^0.8.20;

contract Ballot {
    struct Voter {
        uint weight;
        uint vote;
        bool voted;
        address delegate;
    }

    struct Proposal {
        bytes32 name;
        bytes32 descriptionHash;
        uint voteCount;
    }

    address public chairperson;
    mapping(address => Voter) public voters;
    Proposal[] public proposals;
    uint public votingDeadline;

    event VoteCast(address indexed voter, uint proposal, uint weight);
    event Delegated(address indexed from, address indexed to);
    event RightToVoteGranted(address indexed voter);

    constructor(
        bytes32[] memory proposalNames,
        byte32[] memory proposaldescptionHash,
        uint durationSeconds
    ) {
        chairperson = msg.sender;
        voters[chairperson].weight = 1;
        votingDeadline = block.timestamp + durationSeconds;
        for (uint i = 0; i < proposalNames.length; i++) {
            proposals.push(
                Proposal({
                    name: proposalNames[i],
                    descptionHash: proposaldescptionHash[i],
                    voteCount: 0
                })
            );
        }
    }

    function givenRightToVote(address voter) external {
        require(block.timestamp < votingDeadline, "Voting period has ended.");
        require(
            msg.sender == chairperson,
            "Only chairperson can give right to vote."
        );
        require(!voters[voter].voted, "The voter already voted.");
        require(voters[voter].weight == 0);
        voters[voter].weight = 1;
        emit RightToVoteGranted(voter);
    }

    function delegate(address to) external {
        require(block.timestamp < votingDeadline, "Voting period has ended.");
        Voter storage sender = voters[msg.sender];
        require(sender.weight != 0, "You have no right to vote.");
        require(!sender.voted, "You already voted.");
        require(to != msg.sender, "Self-delegate is disallowed.");
        while (voters[to].delegate != address(0)) {
            to = voters[to].delegate;
            require(to != msg.sender, "Found loop in delegation.");
        }
        Voter storage delegate_ = voters[to];
        require(delegate_.weight >= 1, "Delegate has no right to vote.");
        sender.voted = true;
        sender.delegate = to;
        if (delegate_.voted) {
            proposals[delegate_.vote].voteCount += sender.weight;
        } else {
            delegate_.weight += sender.weight;
        }
        emit Delegated(msg.sender, to);
    }

    function vote(uint proposal) external {
        require(block.timestamp < votingDeadline, "Voting period has ended.");
        Voter storage sender = voters[msg.sender];
        require(sender.weight != 0, "Has no right to vote.");
        require(!sender.voted, "Already voted");
        sender.voted = true;
        sender.vote = proposal;
        proposals[proposal].voteCount += sender.weight;
        emit VoteCast(msg.sender, proposal, sender.weight);
    }

    function winningProposal() public view returns (uint winningProposal_) {
        require(block.timestamp >= votingDeadline, "Voting is still ongoing.");
        uint winningVoteCount = 0;
        for (uint p = 0; p < proposals.length; p++) {
            if (proposals[p].voteCount > winningVoteCount) {
                winningVoteCount = proposals[p].voteCount;
                winningProposal_ = p;
            }
        }
    }

    function winnerName() external view returns (bytes32 winnerName_) {
        require(block.timestamp >= votingDeadline, "Voting is still ongoing.");
        winnerName_ = proposals[winningProposal()].name;
    }
}
