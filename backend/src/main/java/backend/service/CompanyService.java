package backend.service;

import backend.entity.CompanyEntity;
import backend.repository.CompanyRepository;
import com.aggregator.internship.CompanyInternship;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class CompanyService {

    private final CompanyRepository companyRepository;


    @Transactional
    public CompanyEntity findOrCreateCompany(CompanyInternship proto) {
        companyRepository.insertIfNotExists(
                UUID.randomUUID(),
                proto.getCompanyName().toUpperCase(),
                proto.getSourceUrl(),
                proto.getSourceSite()
        );

        return companyRepository.findByCompanyName(proto.getCompanyName().toUpperCase())
                .orElseThrow(() ->
                        new RuntimeException("Company with name " + proto.getCompanyName() + " not found"));
    }


    @Transactional
    public void saveCompaniesBatch(Iterable<CompanyInternship> companies) {
        for (CompanyInternship c : companies) {
            insertIfNotExists(c);
        }
    }

    private void insertIfNotExists(CompanyInternship proto) {
        companyRepository.insertIfNotExists(
                UUID.randomUUID(),
                proto.getCompanyName().toUpperCase(),
                proto.getSourceUrl(),
                proto.getSourceSite()
        );
    }
}