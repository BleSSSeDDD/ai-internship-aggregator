package backend.repository;

import backend.entity.CompanyEntity;
import backend.entity.InternshipEntity;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface InternshipRepository extends JpaRepository<InternshipEntity, UUID> {

    // ИСПРАВЛЕНИЕ: Убрали AndCompanySourceUrl (избыточный JOIN)
    Optional<InternshipEntity> findByCompanyAndPositionName(
            CompanyEntity company,
            String positionName
    );

    Page<InternshipEntity> findAll(Pageable pageable);
}