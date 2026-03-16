package backend.entity;

import jakarta.persistence.*;
import lombok.*;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Entity
@Table(name = "company")
@AllArgsConstructor
@NoArgsConstructor
@Setter
@Getter
@EqualsAndHashCode(onlyExplicitlyIncluded = true)
@ToString(exclude = "internships")
public class CompanyEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @EqualsAndHashCode.Include
    private UUID id;

    @Column(name = "company_name", unique = true, nullable = false)
    private String companyName;

    @Column(name = "source_url")
    private String sourceUrl;

    @Column(name = "source_site")
    private String sourceSite;

    @OneToMany(mappedBy = "company", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<InternshipEntity> internships = new ArrayList<>();

    public void addInternship(InternshipEntity internship) {
        this.internships.add(internship);
        internship.setCompany(this);
    }
}