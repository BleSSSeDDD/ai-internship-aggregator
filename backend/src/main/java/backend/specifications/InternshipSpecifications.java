package backend.specifications;

import backend.entity.CompanyEntity;
import backend.entity.InternshipEntity;
import jakarta.persistence.criteria.Join;
import jakarta.persistence.criteria.JoinType;
import jakarta.persistence.criteria.Predicate;
import org.springframework.data.jpa.domain.Specification;

import java.util.Set;

public class InternshipSpecifications {

    public static Specification<InternshipEntity> hasLocation(String location) {
        if (location == null || location.isBlank()) {
            return null;
        }



        return (root, query, cb) ->
                cb.like(
                        root.get("location"),
                        "%" + location.trim() + "%"
                );
    }

    public static Specification<InternshipEntity> hasCompanyName(String companyName) {
        if (companyName == null || companyName.isBlank()) {
            return null;
        }

        return (root, query, cb) -> {

            Join<InternshipEntity, CompanyEntity> join = root.join("company", JoinType.INNER);

            return cb.like(
                    join.get("companyName"),
                    "%" + companyName.trim().toUpperCase() + "%"
            );
        };
    }

    public static Specification<InternshipEntity> hasMinSalary(Integer salary) {
        if (salary == null) {
            return null;
        }

        return (root, query, cb) ->
                cb.greaterThanOrEqualTo(root.get("minSalary"), salary);
    }

    public static Specification<InternshipEntity> hasTechnologys(Set<String> technologies) {
        if (technologies == null || technologies.isEmpty()) {
            return null;
        }

        return (root, query, cb) -> {
            Predicate[] predicates = technologies.stream()
                    .filter(tech -> tech != null && !tech.isBlank())
                    .map(String::trim)
                    .map(tech -> cb.isMember(tech, root.get("techStack")))
                    .toArray(Predicate[]::new);

            return cb.and(predicates);
        };
    }
}