package backend.dto;

import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

public record InternshipDTO(
        UUID id,
        String positionName,
        String companyName,
        Set<String> techStack,
        Integer minSalary,
        String location,
        String internshipDates,
        String selectionProcess,
        String description,
        String applicationDeadline,
        String contactInfo,
        String experienceRequirements
) {}
