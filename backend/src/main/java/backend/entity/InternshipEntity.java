package backend.entity;

import jakarta.persistence.*;
import lombok.*;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Entity
@Table(name = "internships", indexes = {
        @Index(name = "idx_internship_search", columnList = "company_id, position_name")
}, uniqueConstraints = {
        @UniqueConstraint(columnNames = {"company_id", "position_name"})
})
@AllArgsConstructor
@NoArgsConstructor
@Setter
@Getter
@EqualsAndHashCode(onlyExplicitlyIncluded = true)
@ToString(exclude = "company")
public class InternshipEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @EqualsAndHashCode.Include
    private UUID id;

    @Column(name = "position_name", nullable = false)
    private String positionName;

    @ElementCollection
    @CollectionTable(name = "internship_tech_stack", joinColumns = @JoinColumn(name = "internship_id"))
    @Column(name = "technology")
    private Set<String> techStack = new HashSet<>();

    @Column(name = "min_salary")
    private Integer minSalary;

    @Column(name = "location")
    private String location;

    @Column(name = "internship_dates")
    private String internshipDates;

    @Column(name = "selection_process")
    private String selectionProcess;

    @Column(name = "description", columnDefinition = "text")
    private String description;

    @Column(name = "application_deadline")
    private String applicationDeadline;

    @Column(name = "contact_info")
    private String contactInfo;

    @Column(name = "experience_requirements", columnDefinition = "text")
    private String experienceRequirements;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "company_id", nullable = false)
    private CompanyEntity company;
}