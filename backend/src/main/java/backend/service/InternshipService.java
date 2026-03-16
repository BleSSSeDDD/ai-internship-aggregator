package backend.service;

import backend.dto.InternshipDTO;
import backend.entity.CompanyEntity;
import backend.entity.InternshipEntity;
import backend.mapper.InternshipMapper;
import backend.repository.InternshipRepository;
import com.aggregator.internship.CompanyInternship;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.data.domain.Page;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.PageRequest;

import java.util.*;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class InternshipService {

    private final InternshipRepository internshipRepository;
    private final CompanyService companyService;
    private final InternshipMapper internshipMapper;

    @Transactional
    public InternshipEntity findOrCreateInternship(CompanyInternship companyInternship) {
        CompanyEntity company = companyService.findOrCreateCompany(companyInternship);

        return internshipRepository.findByCompanyAndPositionName(company, companyInternship.getPositionName())
                .orElseGet(() -> saveNewInternship(companyInternship, company));
    }

    private InternshipEntity saveNewInternship(CompanyInternship ci, CompanyEntity company) {
        InternshipEntity internship = buildInternship(ci, company);

        try {
            return internshipRepository.saveAndFlush(internship);
        } catch (DataIntegrityViolationException e) {
            log.warn("Conflict saving internship {}, fetching existing", ci.getPositionName());
            return internshipRepository.findByCompanyAndPositionName(company, ci.getPositionName())
                    .orElseThrow(() -> new RuntimeException("Conflict resolution failed for " + ci.getPositionName()));
        }
    }

    private InternshipEntity buildInternship(CompanyInternship ci, CompanyEntity company) {
        InternshipEntity internship = new InternshipEntity();
        internship.setPositionName(ci.getPositionName());
        internship.setTechStack(new HashSet<>(ci.getTechStackList()));
        internship.setMinSalary(ci.getMinSalary());
        internship.setLocation(ci.getLocation());
        internship.setInternshipDates(ci.getInternshipDates());
        internship.setSelectionProcess(ci.getSelectionProcess());
        internship.setDescription(ci.getDescription());
        internship.setApplicationDeadline(ci.getApplicationDeadline());
        internship.setContactInfo(ci.getContactInfo());
        internship.setExperienceRequirements(ci.getExperienceRequirements());

        company.addInternship(internship);
        return internship;
    }


    @Transactional
    public List<InternshipEntity> saveBatch(List<CompanyInternship> batch) {

        Map<String, CompanyEntity> companies = batch.stream()
                .map(CompanyInternship::getCompanyName)
                .distinct()
                .collect(Collectors.toMap(
                        name -> name,
                        name -> companyService.findOrCreateCompany(
                                batch.stream().filter(c -> c.getCompanyName().equals(name)).findFirst().get()
                        )
                ));

        List<InternshipEntity> internships = new ArrayList<>();
        for (CompanyInternship ci : batch) {
            CompanyEntity company = companies.get(ci.getCompanyName());
            internshipRepository.findByCompanyAndPositionName(company, ci.getPositionName())
                    .ifPresentOrElse(
                            internships::add,
                            () -> internships.add(buildInternship(ci, company))
                    );
        }

        return internshipRepository.saveAll(internships);
    }

    public Page<InternshipDTO> getInternships(int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<InternshipEntity> internshipEntityPage;

        internshipEntityPage =  internshipRepository.findAll(pageable);

        return internshipEntityPage.map(internshipMapper::toDTO);
    }
}