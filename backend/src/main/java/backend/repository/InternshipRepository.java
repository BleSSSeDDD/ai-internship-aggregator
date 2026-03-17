package backend.repository;

import backend.entity.CompanyEntity;
import backend.entity.InternshipEntity;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface InternshipRepository extends JpaRepository<InternshipEntity, UUID>,
        JpaSpecificationExecutor<InternshipEntity> {

    Optional<InternshipEntity> findByCompanyAndPositionName(
            CompanyEntity company,
            String positionName
    );

    Page<InternshipEntity> findAll(Pageable pageable);

    @Query(value = """
    SELECT DISTINCT technology from internship_tech_stack order by technology
        """, nativeQuery = true)
    List<String> getAllTech();

}