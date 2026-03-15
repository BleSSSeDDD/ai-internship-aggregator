package backend.repository;

import backend.entity.CompanyEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.Optional;
import java.util.UUID;

public interface CompanyRepository extends JpaRepository<CompanyEntity, UUID> {
    Optional<CompanyEntity> findByCompanyName(String companyName);

    @Modifying
    @Query(value = """
    INSERT INTO company(id, company_name, source_url, source_site)
    VALUES (:id, :companyName, :sourceUrl, :sourceSite)
    ON CONFLICT (company_name) DO NOTHING
    """, nativeQuery = true)
    void insertIfNotExists(
            @Param("id") UUID id,
            @Param("companyName") String companyName,
            @Param("sourceUrl") String sourceUrl,
            @Param("sourceSite") String sourceSite
    );

}
