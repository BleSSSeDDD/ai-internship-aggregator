package backend.service;

import backend.entity.InternshipEntity;
import backend.repository.InternshipRepository;
import com.aggregator.internship.CompanyInternship;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j; // Добавлено для логирования
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.HashSet;
import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j // Добавлено
public class InternshipService {

    private final InternshipRepository internshipRepository;
    private final CompanyService companyService;

    @Transactional
    public InternshipEntity findOrCreateInternship(CompanyInternship companyInternship) {
        var company = companyService.findOrCreateCompany(companyInternship);

        Optional<InternshipEntity> existing = internshipRepository
                .findByCompanyAndPositionName(
                        company,
                        companyInternship.getPositionName()
                );

        if (existing.isPresent()) {
            return existing.get();
        }

        InternshipEntity internship = new InternshipEntity();
        internship.setPositionName(companyInternship.getPositionName());
        internship.setTechStack(new HashSet<>(companyInternship.getTechStackList()));
        internship.setMinSalary(companyInternship.getMinSalary());
        internship.setLocation(companyInternship.getLocation());
        internship.setInternshipDates(companyInternship.getInternshipDates());
        internship.setSelectionProcess(companyInternship.getSelectionProcess());
        internship.setDescription(companyInternship.getDescription());
        internship.setApplicationDeadline(companyInternship.getApplicationDeadline());
        internship.setContactInfo(companyInternship.getContactInfo());
        internship.setExperienceRequirements(companyInternship.getExperienceRequirements());

        company.addInternship(internship);

        try {
            return internshipRepository.saveAndFlush(internship);
        } catch (DataIntegrityViolationException e) {
            log.warn("DataIntegrityViolation for position: {}. Fetching existing record.",
                    companyInternship.getPositionName());

            return internshipRepository
                    .findByCompanyAndPositionName(
                            company,
                            companyInternship.getPositionName()
                    )
                    .orElseThrow(() -> new RuntimeException(
                            "Conflict resolution failed for: " + companyInternship.getPositionName()
                    ));
        }
    }
}